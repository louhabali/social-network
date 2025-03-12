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

type NotifRes struct {
	Id          int    `json:"id"`
	RelatedId   int    `json:"related_id"`
	Type        string `json:"type"`
	RelatedName string `json:"related_name"`
	ReceverId   int    `json:"recever_id"`
	SenderId    *int   `json:"sender_id"`
}

func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userId, _ := r.Context().Value("userID").(int)
	// fmt.Println(userId)
	query := `
		SELECT 
			n.id, 
			n.type, 
			n.related_id,
			n.recever_id,
			CASE 
				WHEN n.type = 'group event' THEN g.title
				ELSE u.first_name || ' ' || u.last_name
			END AS related_name,
			CASE 
				WHEN n.type = 'group event' THEN null
				ELSE n.sender_id
			END AS sender_id
		FROM notifications AS n
		JOIN users As u ON n.sender_id = u.id
		LEFT JOIN groups AS g ON n.related_id = g.id
		WHERE n.recever_id = ? AND n.is_read = 0;
	`
	rows, err := h.db.Query(query, userId)
	// followReqs, err := getFolowReq(h.db, userId)
	if err != nil {
		log.Println("Error getting notifs: ", err)
		http.Error(w, fmt.Sprintf("Error getting notifs: %v", err), http.StatusInternalServerError)
		return
	}

	var notifications []NotifRes
	for rows.Next() {
		// fmt.Println("turn")
		var notif NotifRes
		rows.Scan(&notif.Id, &notif.Type, &notif.RelatedId, &notif.ReceverId, &notif.RelatedName, &notif.SenderId)
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
	fmt.Println("notif ->", notif)
	_, err := db.Exec(query, notif.Type, notif.RelatedId, notif.ReciverId, notif.SenderId)
	if err != nil {
		fmt.Println("err insert notif", err)
		return err
	}

	return nil
}

func (h *Handler) MarkNotifAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var notif NotifRes
	err := json.NewDecoder(r.Body).Decode(&notif)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	_, err = h.db.Exec("UPDATE notifications SET is_read = 1 WHERE id = ?", notif.Id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
