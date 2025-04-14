package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Conversation struct {
	ChatID      int     `json:"chat_id"`
	ChatName    string  `json:"chat_name"`
	LastMsgTime *string `json:"last_message_time"`
	LastMessage *string `json:"last_message"`
	UnreadCount int     `json:"unread_count"`
}

type Message struct {
	ID         int    `json:"id"`
	ChatID     int    `json:"chat_id"`
	SenderID   int    `json:"sender_id"`
	SenderName string `json:"sender_name"`
	IsGroup    bool   `json:"is_group"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}

// CreateChatAfterApproval creates a chat after a follow request is accepted
func CreateChatAfterApproval(db *sql.DB, followReqID int) error {
	var user1, user2 int
	query := `SELECT follower_id, following_id FROM follow_requests WHERE id = ?`
	err := db.QueryRow(query, followReqID).Scan(&user1, &user2)
	if err != nil {
		return err
	}
	_, err = CreatePrivateChat(db, user1, user2)
	return err
}

// GetConversations retrieves the user's conversations from the database
func GetConversations(db *sql.DB, userID int) ([]Conversation, error) {
	query := `
		SELECT 
    	c.id AS chat_id,
    COALESCE(
        c.name, 
        (SELECT first_name || ' ' || last_name 
         FROM users 
         WHERE id = (SELECT user_id FROM chat_members 
                     WHERE chat_id = c.id AND user_id != ? 
                     LIMIT 1))
    ) AS chat_name,
    	(SELECT content FROM messages WHERE chat_id = c.id ORDER BY id DESC LIMIT 1) AS last_message,
    	(SELECT created_at FROM messages WHERE chat_id = c.id ORDER BY id DESC LIMIT 1) AS last_msg_time,
    	(
        	SELECT COUNT(m.id) 
        	FROM messages m
        	JOIN chat_members cm ON m.chat_id = cm.chat_id
        	WHERE m.chat_id = c.id
        	AND cm.user_id = ?
        	AND m.sender_id != ?
        	AND (cm.last_read_message_id IS NULL OR m.id > cm.last_read_message_id)
    	) AS unread_count
		FROM chat_members cm
		JOIN chats c ON cm.chat_id = c.id
		WHERE cm.user_id = ?
		ORDER BY 
    		(SELECT id FROM messages WHERE chat_id = c.id ORDER BY id DESC LIMIT 1) DESC;`

	rows, err := db.Query(query, userID, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation

		err := rows.Scan(
			&conv.ChatID,
			&conv.ChatName,
			&conv.LastMessage,
			&conv.LastMsgTime,
			&conv.UnreadCount,
		)
		if err != nil {
			return nil, err
		}

		conversations = append(conversations, conv)
	}

	return conversations, nil
}

// GetMessages fetches messages for a given chat ID with pagination
func GetMessages(db *sql.DB, userID, chatID, offset int) ([]Message, error) {
	rows, err := db.Query(`SELECT 
		m.id,
		m.chat_id,
		m.sender_id,
		u.first_name || ' ' || u.last_name,
		m.content,
		m.created_at,
		c.is_group
		FROM messages as m 
		left join users as u on u.id = m.sender_id
		left join chats as c on c.id = m.chat_id
		WHERE m.chat_id = ? ORDER BY m.created_at DESC LIMIT 10 OFFSET ?`, chatID, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.SenderName, &msg.Content, &msg.CreatedAt, &msg.IsGroup); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	MarkMessagesAsRead(db, chatID, userID)
	return messages, nil
}

// InsertMessage inserts a new message into the database
func InsertMessage(db *sql.DB, senderID, chatID int, content string) (Message, error) {
	if !CanInsert(db, chatID, senderID) {
		return Message{}, fmt.Errorf("At least one follow the other constraint")
	}

	time := time.Now().Format("2006-01-02 15:04:05")

	result, err := db.Exec(`INSERT INTO messages (chat_id, sender_id, content) VALUES (?, ?, ?)`, chatID, senderID, content)
	if err != nil {
		return Message{}, err
	}

	msgID, _ := result.LastInsertId()
	msg := Message{
		ID:        int(msgID),
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   content,
		CreatedAt: time,
	}

	_, err = db.Exec(
		"UPDATE chat_members SET last_read_message_id = ? WHERE chat_id = ? AND user_id = ?",
		msgID, chatID, senderID,
	)
	if err != nil {
		return Message{}, err
	}

	return msg, nil
}

// GetUsersInChat retrieves the user IDs of members in a chat
func GetUsersInChat(db *sql.DB, chatID int) ([]int, error) {
	rows, err := db.Query("SELECT user_id FROM chat_members WHERE chat_id = ?", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, nil
}

// CreatePrivateChat creates a private chat between two users
func CreatePrivateChat(db *sql.DB, user1, user2 int) (int, error) {
	var chatID int
	err := db.QueryRow(`SELECT chat_id FROM chat_members WHERE user_id IN (?, ?) GROUP BY chat_id HAVING COUNT(DISTINCT user_id) = 2`, user1, user2).Scan(&chatID)

	if err == nil {
		return chatID, nil
	} else if err != sql.ErrNoRows {
		return 0, err
	}

	res, err := db.Exec("INSERT INTO chats (name, is_group) VALUES (NULL, 0)")
	if err != nil {
		return 0, err
	}

	newChatID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, userID := range []int{user1, user2} {
		if err := AddChatMember(db, int(newChatID), userID); err != nil {
			return 0, err
		}
	}
	return int(newChatID), nil
}

// AddChatMember adds a user to a chat
func AddChatMember(db *sql.DB, chatID, userID int) error {
	_, err := db.Exec("INSERT INTO chat_members (chat_id, user_id) VALUES (?, ?)", chatID, userID)
	return err
}

func CreatGroupChat(db *sql.DB, adminId int, name string) (int, error) {
	res, err := db.Exec("INSERT INTO chats (name, is_group) VALUES (?, 1)", name)
	if err != nil {
		return 0, err
	}

	chatID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = AddChatMember(db, int(chatID), adminId)
	if err != nil {
		return 0, err
	}

	return int(chatID), nil
}

func MarkMessagesAsRead(db *sql.DB, chat_id, userID int) error {
	var latestMessageID int64
	err := db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM messages WHERE chat_id = ?", chat_id).Scan(&latestMessageID)
	if err != nil {
		return err
	}

	if latestMessageID == 0 {
		return nil
	}

	_, err = db.Exec(
		"UPDATE chat_members SET last_read_message_id = ? WHERE chat_id = ? AND user_id = ?", latestMessageID, chat_id, userID,
	)
	return err
}

func GetChatIDByGroupName(db *sql.DB, groupName string) (int, error) {
	var chatID int
	err := db.QueryRow("SELECT id FROM chats WHERE name = ?", groupName).Scan(&chatID)
	if err != nil {
		return 0, err
	}
	return chatID, nil
}

func CanInsert(db *sql.DB, chatId, senderId int) bool {
	query := `
		SELECT 1
		FROM chats c
		JOIN chat_members cm1 ON c.id = cm1.chat_id
		WHERE c.id = ?
		AND cm1.user_id = ?
		AND (
		    c.is_group = 1
		    OR EXISTS (
		        SELECT 1
		        FROM follow_requests f
		        JOIN chat_members cm2 ON cm2.chat_id = c.id
		        WHERE f.status = 'approved'
		        AND cm2.user_id != cm1.user_id
		        AND (
		            (f.follower_id = cm1.user_id AND f.following_id = cm2.user_id)
		            OR
		            (f.follower_id = cm2.user_id AND f.following_id = cm1.user_id)
		        )
		    )
		);
	`

	var exists int
	err := db.QueryRow(query, chatId, senderId).Scan(&exists)
	if err != nil {
		return false
	}

	return true
}
