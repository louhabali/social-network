package handlers

import (
	"encoding/json"
	"main/pkg/app/models"
	"net/http"
)

func (h *Handler) GetAllusers(w http.ResponseWriter, r *http.Request) {
	myId := r.Context().Value("userID").(int)

	followers, err := models.GetFollowing(h.db, myId)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}
