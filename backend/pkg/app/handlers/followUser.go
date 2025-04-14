package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"main/pkg/app/models"
)

func (h *Handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	followerID := r.Context().Value("userID").(int)

	if r.Method == "POST" {

		isPublic, err := models.CheckIfPublic(h.db, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		status := "pending"
		if isPublic {
			status = "approved"
		}

		reqID, err := models.InsertFollowRequest(h.db, followerID, userID, status)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if isPublic {
			_, err = models.CreatePrivateChat(h.db, followerID, userID)
			if err != nil {
				http.Error(w, "Error creating chat", http.StatusInternalServerError)
				return
			}
			user, _ := models.GetUserInfo(h.db, followerID)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
		} else {
			notif, err := models.InsertNotification(h.db, &models.Notification{
				SenderId:  followerID,
				ReceverId: userID,
				RelatedId: int(reqID),
				Type:      "follow request",
			})
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			wsManager.SendToUsers([]int{notif.ReceverId}, WebSocketMessage{
				Type:    "notification",
				Payload: notif,
			})
		}
	} else if r.Method == "DELETE" {
		err := models.DeleteFollowRequest(h.db, followerID, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		err = models.DeleteNotification(h.db, 0, followerID, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) AcceptFollowReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var notif models.Notification
	if err := json.NewDecoder(r.Body).Decode(&notif); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := models.ApproveFollowRequest(h.db, notif.RelatedId); err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	if err := models.DeleteNotification(h.db, notif.Id, 0, 0); err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	if err := models.CreateChatAfterApproval(h.db, notif.RelatedId); err != nil {
		http.Error(w, "Error creating chat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) RejectFollowReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var notif models.Notification
	if err := json.NewDecoder(r.Body).Decode(&notif); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println(notif)
	if err := models.DeleteFollowRequest(h.db, notif.SenderId, notif.ReceverId); err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	if err := models.DeleteNotification(h.db, notif.Id, 0, 0); err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
