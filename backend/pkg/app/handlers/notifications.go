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
	Id         int     `json:"id"`
	RelatedId  int     `json:"related_id"`
	Type       string  `json:"type"`
	SenderName string  `json:"sender_name"`
	SenderId   int     `json:"sender_id"`
	ReceverId  int     `json:"recever_id"`
	GroupId    *int    `json:"group_id"`
	GroupName  *string `json:"group_name"`
}

func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	userId, _ := r.Context().Value("userID").(int)
	query := `
		SELECT 
		    n.id, 
		    n.type, 
		    n.related_id,
		    n.sender_id,
			n.recever_id,
		    COALESCE(ge.group_id, gm.group_id) AS group_id,
		    u.first_name || ' ' || u.last_name AS sender_name,
		    g.title AS group_name
		FROM notifications AS n
		JOIN users As u ON n.sender_id = u.id
		LEFT JOIN group_events AS ge ON ((n.type = 'group event') AND n.related_id = ge.id)
		LEFT JOIN group_members AS gm ON  ((n.type = 'group join invatation' OR n.type = 'group join request') AND n.related_id = gm.id)
		LEFT JOIN groups AS g ON g.id = COALESCE(ge.group_id, gm.group_id)
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
		rows.Scan(&notif.Id, &notif.Type, &notif.RelatedId, &notif.SenderId, &notif.ReceverId, &notif.GroupId, &notif.SenderName, &notif.GroupName)
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

func (h *Handler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	userId, _ := r.Context().Value("userID").(int)
	var notifsNum int
	err := h.db.QueryRow("SELECT COUNT(*) FROM notifications WHERE recever_id = ? and is_read = 0;", userId).Scan(&notifsNum)
	if err != nil {
		fmt.Println("notifs numb ->", err)
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]int{"unread_count": notifsNum}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}
