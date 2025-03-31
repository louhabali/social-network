package models

import (
	"database/sql"
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
	ID        int    `json:"id"`
	ChatID    int    `json:"chat_id"`
	SenderID  int    `json:"sender_id"`
	IsSender  bool   `json:"is_sender"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
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
		WITH last_messages AS (
		    SELECT 
		        chat_id, 
		        id AS last_message_id, 
		        content AS last_message, 
		        created_at AS last_message_time
		    FROM messages
		    WHERE id IN (
		        SELECT MAX(id) FROM messages GROUP BY chat_id
		    )
		),
		unread_counts AS (
		    SELECT 
		        m.chat_id, 
		        COUNT(m.id) AS unread_count
		    FROM messages m
		    LEFT JOIN user_lastMsg_read ulr 
		        ON m.chat_id = ulr.chat_id AND ulr.user_id = ?
		    WHERE m.id > COALESCE((SELECT last_read_message_id 
		                           FROM user_lastMsg_read 
		                           WHERE user_id = ? AND chat_id = m.chat_id), 0)
		    GROUP BY m.chat_id
		)
		SELECT 
		    c.id AS chat_id,
		    COALESCE(c.name, u.first_name || ' ' || u.last_name) AS chat_name,
		    lm.last_message, 
		    lm.last_message_time, 
		    COALESCE(uc.unread_count, 0) AS unread_count
		FROM chat_members cm
		JOIN chats c ON cm.chat_id = c.id
		LEFT JOIN last_messages lm ON c.id = lm.chat_id
		LEFT JOIN unread_counts uc ON c.id = uc.chat_id
		LEFT JOIN chat_members other_cm ON c.id = other_cm.chat_id AND other_cm.user_id != ?
		LEFT JOIN users u ON other_cm.user_id = u.id
		WHERE cm.user_id = ?
	`
	rows, err := db.Query(query, userID, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		err := rows.Scan(&conv.ChatID, &conv.ChatName, &conv.LastMessage, &conv.LastMsgTime, &conv.UnreadCount)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}
	return conversations, nil
}

// GetMessages fetches messages for a given chat ID with pagination
func GetMessages(db *sql.DB, userID, chatID, offset int) ([]Message, error) {
	rows, err := db.Query(`SELECT id, chat_id, sender_id, content, created_at FROM messages WHERE chat_id = ? ORDER BY created_at DESC LIMIT 10 OFFSET ?`, chatID, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		msg.IsSender = msg.SenderID == userID
		messages = append(messages, msg)
	}
	return messages, nil
}

// InsertMessage inserts a new message into the database
func InsertMessage(db *sql.DB, senderID, chatID int, content string) (Message, error) {
	time := time.Now().Format("2006-01-02 15:04:05")

	result, err := db.Exec( `INSERT INTO messages (chat_id, sender_id, content, created_at) VALUES (?, ?, ?, ?)`, chatID, senderID, content, time)
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

	return msg, nil
}

// UpdateLastReadMessage updates or inserts the last read message for a user in a chat
func UpdateLastReadMessage(db *sql.DB, userID, chatID int) error {
	var lastMessageID int
	err := db.QueryRow(`SELECT id FROM messages WHERE chat_id = ? ORDER BY created_at DESC LIMIT 1`, chatID).Scan(&lastMessageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	var exists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM user_lastMsg_read WHERE chat_id = ? AND user_id = ?)`, chatID, userID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = db.Exec(`UPDATE user_lastMsg_read SET last_read_message_id = ? WHERE chat_id = ? AND user_id = ?`, lastMessageID, chatID, userID)
	} else {
		_, err = db.Exec(`INSERT INTO user_lastMsg_read (user_id, chat_id, last_read_message_id) VALUES (?, ?, ?)`, userID, chatID, lastMessageID)
	}
	return err
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

// CreateGroupChat creates a new group chat
func CreateGroupChat(db *sql.DB, adminID int, name string) (int, error) {
	res, err := db.Exec("INSERT INTO chats (name, is_group) VALUES (?, 1)", name)
	if err != nil {
		return 0, err
	}

	chatID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = AddChatMember(db, int(chatID), adminID)
	if err != nil {
		return 0, err
	}
	return int(chatID), nil
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