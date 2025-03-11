package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type notificationToInsert struct {
	Type      string
	RelatedId int
	ReciverId int
	SenderId  int
}

func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userId, _ := r.Context().Value("userID").(int)
	// fmt.Println(userId)
	query := `
		SELECT 
		    n.id, 
		    n.type, 
		    n.related_id, 
		    CASE 
		        WHEN n.type = 'group event' THEN g.name 
		        ELSE u.first_name || ' ' || u.last_name 
		    END AS related_name
		FROM notifications AS n
		JOIN users AS u ON n.sender_id = u.id
		LEFT JOIN groups AS g ON n.type = 'group event' AND n.related_id = g.id
		WHERE n.reciver_id = ? AND n.is_read = 0
	`
	rows, err := h.db.Query(query, userId)
	// followReqs, err := getFolowReq(h.db, userId)
	if err != nil {
		log.Println("Error getting notifs: ", err)
		http.Error(w, fmt.Sprintf("Error getting notifs: %v", err), http.StatusInternalServerError)
		return
	}

	type notifRes struct {
		Id          int    `json:"id"`
		RelatedId   int    `json:"related_id"`
		Type        string `json:"type"`
		RelatedName string `json:"related_name"`
	}

	var notifications []notifRes
	for rows.Next() {
		fmt.Println("turn")
		var notif notifRes
		rows.Scan(&notif.Id, &notif.Type, &notif.RelatedId, &notif.RelatedName)
		notifications = append(notifications, notif)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(notifications); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func insertNotification(db *sql.DB, notif *notificationToInsert) error {
	query := `
		INSERT INTO notifications (type, related_id, recever_id, sender_id, is_read, created_at)
		VALUES (?, ?, ?, ?, 0, CURRENT_TIMESTAMP)
	`
	fmt.Println(notif)
	_, err := db.Exec(query, notif.Type, notif.RelatedId, notif.ReciverId, notif.SenderId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
