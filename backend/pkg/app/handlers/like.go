package handlers

import (
	"encoding/json"
	"net/http"

	"main/pkg/app/models"
)

func (h *Handler) LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var data struct {
		Pid      int `json:"postId"`
		Interact int `json:"likeStatus"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	userID := r.Context().Value("userID").(int)

	err = models.InsertOrUpdateLike(h.db, userID, data.Pid, data.Interact)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
