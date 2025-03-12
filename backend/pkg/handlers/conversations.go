package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Conversation struct {
	ChatID      *int     `json:"chat_id"`
	UserID      *int     `json:"user_id"`
	ChatName    *string  `json:"chat_name"`
	LastMsgTime *string  `json:"last_message_time"`
	LastMessage *string `json:"last_message"`
	UnreadCount *int     `json:"unread_count"`
}



func (h *Handler) GetConversations(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(int)
	
	query := `

	`
	rows, err := h.db.Query(query, userID, userID, userID, userID, userID)
	if err != nil {
		fmt.Println("err in query", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		fmt.Println("turns")
		var conv Conversation
		err := rows.Scan(&conv.ChatID, &conv.ChatName, &conv.LastMessage, &conv.UnreadCount, &conv.LastMsgTime, &conv.UserID)
		if err != nil {
			fmt.Println("err in scan =>", err)
			http.Error(w, "Error reading data", http.StatusInternalServerError)
			return
		}
		conversations = append(conversations, conv)
	}

	fmt.Println(conversations)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}
