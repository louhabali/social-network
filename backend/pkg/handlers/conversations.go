package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Conversation struct {
	ChatID      int     `json:"chat_id"`
	ChatName    string  `json:"chat_name"`
	LastMsgTime *string  `json:"last_message_time"`
	LastMessage *string `json:"last_message"`
	UnreadCount int     `json:"unread_count"`
}



func (h *Handler) GetConversations(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(int)
	
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
	rows, err := h.db.Query(query, userID, userID, userID, userID)
	if err != nil {
		fmt.Println("err in query", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		// fmt.Println("turns")
		var conv Conversation
		err := rows.Scan(&conv.ChatID, &conv.ChatName, &conv.LastMessage, &conv.LastMsgTime, &conv.UnreadCount)
		if err != nil {
			fmt.Println("err in scan =>", err)
			http.Error(w, "Error reading data", http.StatusInternalServerError)
			return
		}
		conversations = append(conversations, conv)
	}

	// fmt.Println(conversations)
	w.Header().Set("Content-Type", "application/json")
	
	err = json.NewEncoder(w).Encode(conversations)
	if err != nil {
		// fmt.Println("err in scan =>", err)
		http.Error(w, "Error in encodin json", http.StatusInternalServerError)
		return
	}
}
