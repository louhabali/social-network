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

/*
to select just folowing and follower
potential_chats AS (
    SELECT
        NULL AS chat_id,
        u.nickname AS chat_name,
        NULL AS last_message,
        NULL AS unread_count,
        NULL AS last_message_time,
        u.id AS user_id
    FROM users u
    WHERE u.id IN (
        SELECT following_id FROM follow_requests WHERE follower_id = ? AND status = 'accepted'
        UNION
        SELECT follower_id FROM follow_requests WHERE following_id = ? AND status = 'accepted'
    )
    AND u.id != ?
    AND NOT EXISTS (
        SELECT 1 FROM chat_members cm
        JOIN chats c ON cm.chat_id = c.id
        WHERE cm.user_id = u.id AND cm.chat_id IN (SELECT id FROM chats)
    )
)
*/

func (h *Handler) GetConversations(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(int)
	
	query := `
		WITH existing_chats AS (
		    SELECT 
		        c.id AS chat_id,
		        CASE 
		            WHEN c.is_group = 1 THEN c.name 
		            ELSE (
		                SELECT u.nickname FROM users u 
		                WHERE u.id = (SELECT cm.user_id FROM chat_members cm 
		                              WHERE cm.chat_id = c.id AND cm.user_id != ? LIMIT 1)
		            )
		        END AS chat_name,
		        (SELECT content FROM messages m 
		         WHERE m.chat_id = c.id ORDER BY m.created_at DESC LIMIT 1) AS last_message,
				COALESCE(
				    (SELECT COUNT(*) FROM messages m 
				     WHERE m.chat_id = c.id AND m.id > 
				         COALESCE((SELECT last_read_message_id 
				                   FROM user_lastMsg_read ulr 
				                   WHERE ulr.user_id = ? AND ulr.chat_id = c.id), 0)
				    ), 
				    (SELECT COUNT(*) FROM messages m WHERE m.chat_id = c.id) -- If no entry in last read table, count all messages
				) AS unread_count,

		        (SELECT created_at FROM messages m 
		         WHERE m.chat_id = c.id ORDER BY m.created_at DESC LIMIT 1) AS last_message_time,
		        NULL AS user_id
		    FROM chats c
		    JOIN chat_members cm ON c.id = cm.chat_id
		    WHERE cm.user_id = ?
		),
		potential_chats AS (
		    SELECT 
		        NULL AS chat_id,
		        u.nickname AS chat_name,
		        NULL AS last_message,
		        NULL AS unread_count,
		        NULL AS last_message_time,
		        u.id AS user_id
		    FROM users u
		    WHERE u.id != ?
		    AND NOT EXISTS (
		        SELECT 1 FROM chat_members cm 
		        JOIN chats c ON cm.chat_id = c.id 
		        WHERE cm.user_id = u.id 
		        AND cm.chat_id IN (
		            SELECT chat_id FROM chat_members WHERE user_id = ?
		        )
			)
		)
		SELECT * FROM existing_chats
		UNION ALL
		SELECT * FROM potential_chats;
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
