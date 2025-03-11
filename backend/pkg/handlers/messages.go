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
	ID         int    `json:"id"`
	ChatID     int    `json:"chat_id"`
	SenderID   int    `json:"sender_id"`
	IsSender bool   `json:"is_sender"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(int)
	chatID := r.URL.Query().Get("chat_id")
	fmt.Println(chatID)
	if chatID == "" {
		http.Error(w, "chat_id is required", http.StatusBadRequest)
		return
	}

	limit := 10
	offset := 0
	if r.URL.Query().Get("offset") != "" {
		offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	}

	rows, err := h.db.Query(`
	SELECT id, chat_id, sender_id, content, created_at 
	FROM messages 
	WHERE chat_id = ? 
	ORDER BY created_at DESC 
	LIMIT ? OFFSET ?`, chatID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.CreatedAt); err != nil {
			http.Error(w, "Error scanning messages", http.StatusInternalServerError)
			return
		}
		msg.IsSender = msg.SenderID == userID
		messages = append(messages, msg)
	}
	fmt.Println(messages)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func InsertMessage(sendeId, chatId int, content string, db *sql.DB) (Message, error) {
	// var newMsg struct {
	// 	ChatID   int    `json:"chat_id"`
	// 	SenderID int    `json:"sender_id"`
	// 	Content  string `json:"content"`
	// }

	// if err := json.NewDecoder(r.Body).Decode(&newMsg); err != nil {
	// 	http.Error(w, "Invalid JSON", http.StatusBadRequest)
	// 	return
	// }

	query := `INSERT INTO messages (chat_id, sender_id, content, created_at) VALUES (?, ?, ?, ?)`
	time := time.Now().Format("2006-01-02 15:04:05")
	result, err := db.Exec(query, chatId, sendeId, content, time)
	if err != nil {
		// http.Error(w, "Failed to insert message", http.StatusInternalServerError)
		return Message{}, err
	}

	// messageID, _ := result.LastInsertId()
	// msgResponse := Message{
	// 	ID:        int(messageID),
	// 	ChatID:    newMsg.ChatID,
	// 	SenderID:  newMsg.SenderID,
	// 	Content:   newMsg.Content,
	// 	CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	// }

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(msgResponse)
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
	var ulmr struct {
		ChatID            int `json:"chat_id"`
		UserID            int `json:"user_id"`
		LastReadMessageID int `json:"last_read_message_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&ulmr); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var exists bool
	err := h.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM user_lastMsg_read WHERE chat_id = ? AND user_id = ?)`, ulmr.ChatID, ulmr.UserID).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if exists {
		_, err = h.db.Exec(`UPDATE user_lastMsg_read SET last_read_message_id = ? WHERE chat_id = ? AND user_id = ?`, ulmr.LastReadMessageID, ulmr.ChatID, ulmr.UserID)
		if err != nil {
			http.Error(w, "Failed to update last read message", http.StatusInternalServerError)
			return
		}
	} else {
		_, err = h.db.Exec(`INSERT INTO user_lastMsg_read (user_id, chat_id, last_read_message_id) VALUES (?, ?, ?)`, ulmr.UserID, ulmr.ChatID, ulmr.LastReadMessageID)
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
	res, err := db.Exec("INSERT INTO chats (name, is_group) VALUES (NULL, 0)")
	if err != nil {
		return 0, err
	}

	chatID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	chatUsers := []int{user1, user2}
	for _, userId := range chatUsers {
		fmt.Println(userId)
		err := addChatMember(db, int(chatID), userId)
		if err != nil {
			fmt.Println("fail to register chat member in private chat")
			return 0, err
		}
	}

	return int(chatID), nil
}

func addChatMember(db *sql.DB, chatID int, userID int) error {
	_, err := db.Exec("INSERT INTO chat_members (chat_id, user_id) VALUES (?, ?)", chatID, userID)
	return err
}
