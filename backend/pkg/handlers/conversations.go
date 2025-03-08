package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Conversation struct {
	ChatID      int     `json:"chat_id"`
	ChatName    string  `json:"chat_name"`
	LastMessage *string `json:"last_message"`
	UnreadCount int     `json:"unread_count"`
}


func (h *Handler) GetConversations(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(int)
	query := `
		SELECT 
			c.id AS chat_id,
			CASE 
				WHEN c.is_group = 1 THEN c.name 
				ELSE (SELECT u.nickname FROM users u 
					WHERE u.id = (SELECT cm.user_id FROM chat_members cm 
									WHERE cm.chat_id = c.id AND cm.user_id != ? LIMIT 1))
			END AS chat_name,
			(SELECT content FROM messages m 
			WHERE m.chat_id = c.id ORDER BY m.created_at DESC LIMIT 1) AS last_message,
			(SELECT COUNT(*) FROM messages m 
			WHERE m.chat_id = c.id AND m.id > 
				(SELECT COALESCE(last_read_message_id, 0) FROM user_lastMsg_read ulr 
				WHERE ulr.user_id = ? AND ulr.chat_id = c.id)
			) AS unread_count
		FROM chats c
		JOIN chat_members cm ON c.id = cm.chat_id
		WHERE cm.user_id = ?
		ORDER BY 
			(SELECT created_at FROM messages m 
			WHERE m.chat_id = c.id ORDER BY m.created_at DESC LIMIT 1) DESC;
	`
	rows, err := h.db.Query(query, userID, userID, userID)
	if err != nil {
		fmt.Println("err")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	var conversations []Conversation
	for rows.Next() {
		var convo Conversation
		err := rows.Scan(&convo.ChatID, &convo.ChatName, &convo.LastMessage, &convo.UnreadCount)
		if err != nil {
			fmt.Println("err")
			http.Error(w, "Error reading data", http.StatusInternalServerError)
			return
		}
		conversations = append(conversations, convo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
	fmt.Println("req done", conversations)
}
