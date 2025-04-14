package handlers

import (
	"encoding/json"
	"net/http"

	"main/pkg/app/models"
)

func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	userId := r.Context().Value("userID").(int)

	notifications, err := models.GetNotifications(h.db, userId)
	if err != nil {
		http.Error(w, "Error getting notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func (h *Handler) MarkNotifAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	userId := r.Context().Value("userID").(int)
	// var notif models.Notification
	// err := json.NewDecoder(r.Body).Decode(&notif)
	// if err != nil {
	// 	http.Error(w, "Invalid request body", http.StatusBadRequest)
	// 	return
	// }
	// r.Body.Close()

	err := models.MarkNotificationAsRead(h.db, userId)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userId := r.Context().Value("userID").(int)

	unreadCount, err := models.GetUnreadCount(h.db, userId)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	response := map[string]int{"unread_count": unreadCount}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func (h *Handler) DeleteNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	notif := struct {
		Id int `json:"id"` 
	}{}

	err := json.NewDecoder(r.Body).Decode(&notif)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	err = models.DeleteNotification(h.db, notif.Id, 0, 0)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}