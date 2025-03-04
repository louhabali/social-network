package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.PathValue("userId"))
	followerID, ok := r.Context().Value("userID").(int)
	if !ok {
		fmt.Println("invalid userId")
		return
	}

	if r.Method == "POST" {
		_, err := h.db.Exec("INSERT INTO follow_requests (follower_id, following_id) VALUES (?, ?)", followerID, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	} else if r.Method == "DELETE" {
		_, err := h.db.Exec("DELETE FROM follow_requests WHERE follower_id = ? AND following_id = ?", followerID, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
