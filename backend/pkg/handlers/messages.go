package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Message struct {
	ID        int    `json:"id"`
	ChatID    int    `json:"chat_id"`
	SenderID  int    `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chat_id")
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
		if err := rows.Scan(&msg.SenderID, &msg.Content, &msg.CreatedAt); err != nil {
			http.Error(w, "Error scanning messages", http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *Handler) insertMessage(w http.ResponseWriter, r *http.Request) {
	var newMsg struct {
		ChatID   int    `json:"chat_id"`
		SenderID int    `json:"sender_id"`
		Content  string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newMsg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}


	query := `INSERT INTO messages (chat_id, sender_id, content, created_at) VALUES (?, ?, ?, ?)`
	_, err := h.db.Exec(query, newMsg.ChatID, newMsg.SenderID, newMsg.Content, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		http.Error(w, "Failed to insert message", http.StatusInternalServerError)
		return
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
}


func (h *Handler) UpdateLastReadMessage(w http.ResponseWriter, r *http.Request) {
	var UpdateLastReadRequest struct {
		ChatID          int `json:"chat_id"`
		UserID          int `json:"user_id"`
		LastReadMessageID int `json:"last_read_message_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var exists bool
	err := h.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM user_lastMsg_read WHERE chat_id = ? AND user_id = ?)`, req.ChatID, req.UserID).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}