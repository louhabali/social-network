package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Message struct {
	ID        int    `json:"id"`
	ChatID    int    `json:"chat_id"`
	SenderID  int    `json:"sender_id"`
	SenderNm  string `json:"sender_name"`
	IsSender  bool   `json:"is_sender"`
	IsGroup   bool   `json:"is_group"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	userID, _ := r.Context().Value("userID").(int)
	chatID := r.URL.Query().Get("chat_id")
	fmt.Println(chatID)

	if chatID == "" {
		http.Error(w, "chat_id is required", http.StatusBadRequest)
		return
	}

	chatId, _ := strconv.Atoi(chatID)
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	rows, err := h.db.Query(`
	SELECT 
		msgs.id,
		msgs.chat_id,
		msgs.sender_id,
		msgs.content,
		msgs.created_at,
		c.is_group,
		u.first_name || ' ' || u.last_name
	FROM messages as msgs
	LEFT JOIN chats as c ON msgs.chat_id = c.id
	LEFT JOIN users as u ON msgs.sender_id = u.id
	WHERE msgs.chat_id = ?
	ORDER BY msgs.created_at DESC
	LIMIT 10 OFFSET ?`, chatId, offset)
	if err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.CreatedAt, &msg.IsGroup, &msg.SenderNm); err != nil {
			http.Error(w, "Error scanning messages", http.StatusInternalServerError)
			return
		}
		msg.IsSender = msg.SenderID == userID
		messages = append(messages, msg)
	}
	// fmt.Println(messages)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func InsertMessage(sendeId, chatId int, content string, db *sql.DB) (Message, error) {
	// check for private chats if at least one of the member follow the other
	if !canInsert(db, chatId, sendeId) {
		return Message{}, fmt.Errorf("At least one follow the other constraint")
	}

	query := `INSERT INTO messages (chat_id, sender_id, content, created_at) VALUES (?, ?, ?, ?)`
	time := time.Now().Format("2006-01-02 15:04:05")
	result, err := db.Exec(query, chatId, sendeId, content, time)
	if err != nil {
		return Message{}, err
	}

	msgId, _ := result.LastInsertId()
	msg := Message{
		ID:        int(msgId),
		ChatID:    chatId,
		SenderID:  sendeId,
		Content:   content,
		CreatedAt: time,
	}

	return msg, err
}

func (h *Handler) UpdateLastReadMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update last read message")
	if r.Method != http.MethodPost {
		return
	}
	userID, _ := r.Context().Value("userID").(int)

	var ulmr struct {
		ChatID int `json:"chat_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&ulmr); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var lastMessageID int
	err := h.db.QueryRow(`
		SELECT id 
		FROM messages 
		WHERE chat_id = ? 
		ORDER BY created_at DESC LIMIT 1`, ulmr.ChatID).Scan(&lastMessageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}

		http.Error(w, "Database error while fetching last message", http.StatusInternalServerError)
		return
	}

	var exists bool
	err = h.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM user_lastMsg_read WHERE chat_id = ? AND user_id = ?)`, ulmr.ChatID, userID).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error while checking record existence", http.StatusInternalServerError)
		return
	}

	if exists {
		_, err = h.db.Exec(`
			UPDATE user_lastMsg_read 
			SET last_read_message_id = ? 
			WHERE chat_id = ? AND user_id = ?`, lastMessageID, ulmr.ChatID, userID)
		if err != nil {
			http.Error(w, "Failed to update last read message", http.StatusInternalServerError)
			return
		}
	} else {

		_, err = h.db.Exec(`
			INSERT INTO user_lastMsg_read (user_id, chat_id, last_read_message_id) 
			VALUES (?, ?, ?)`, userID, ulmr.ChatID, lastMessageID)
		if err != nil {
			http.Error(w, "Failed to create last read message entry", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func getUsersInChat(chatID int, db *sql.DB) ([]int, error) {
	query := "SELECT user_id FROM chat_members WHERE chat_id = ?"
	rows, err := db.Query(query, chatID)
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

func createPrivateChat(db *sql.DB, user1, user2 int) (int, error) {
	var chatID int
	err := db.QueryRow(`
		SELECT chat_id 
		FROM chat_members 
		WHERE user_id IN (?, ?) 
		GROUP BY chat_id 
		HAVING COUNT(DISTINCT user_id) = 2
	`, user1, user2).Scan(&chatID)

	if err == nil {
		return chatID, nil
	} else if err != sql.ErrNoRows {
		fmt.Println("error insert chat in checkin if it exist", err)
		return 0, err
	}

	res, err := db.Exec("INSERT INTO chats (name, is_group) VALUES (NULL, 0)")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	newChatID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	chatUsers := []int{user1, user2}
	for _, userId := range chatUsers {
		fmt.Println(userId)
		err := addChatMember(db, int(newChatID), userId)
		if err != nil {
			fmt.Println("fail to register chat member in private chat")
			return 0, err
		}
	}

	return int(newChatID), nil
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

	err = addChatMember(db, int(chatID), adminId)
	if err != nil {
		return 0, err
	}

	return int(chatID), nil
}

func addChatMember(db *sql.DB, chatID int, userID int) error {
	_, err := db.Exec("INSERT INTO chat_members (chat_id, user_id) VALUES (?, ?)", chatID, userID)
	return err
}

func canInsert(db *sql.DB, chatId, senderId int) bool {
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
