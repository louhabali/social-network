package handlers

import (
	"encoding/json"
	"fmt"
	"main/pkg/app/models"
	"net/http"
	"strconv"
)

func (h *Handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.PathValue("userId"))
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
			_,err = models.CreatePrivateChat(h.db, followerID, userID)
			if err != nil {
				http.Error(w, "Error creating chat", http.StatusInternalServerError)
				return
			}
			user, _ := models.GetUserInfo(h.db, followerID)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
		} else {
			err = models.InsertNotification(h.db, &models.Notification{
				SenderId:  followerID,
				ReceverId: userID,
				RelatedId: int(reqID),
				Type:      "follow request",
			})
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
		}
	} else if r.Method == "DELETE" {
		err := models.DeleteFollowRequest(h.db, followerID, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		err = models.DeleteNotification(h.db, userID)
		if err != nil {
			fmt.Println(err)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AcceptFollowReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var notif models.Notification
	err := json.NewDecoder(r.Body).Decode(&notif)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	err = models.ApproveFollowRequest(h.db, notif.RelatedId)
	if err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	err = models.DeleteNotification(h.db, notif.Id)
	if err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	err = models.CreateChatAfterApproval(h.db, notif.RelatedId)
	if err != nil {
		fmt.Println(err)
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
	err := json.NewDecoder(r.Body).Decode(&notif)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	err = models.DeleteFollowRequest(h.db, notif.RelatedId, 0)
	if err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	err = models.DeleteNotification(h.db, notif.Id)
	if err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
