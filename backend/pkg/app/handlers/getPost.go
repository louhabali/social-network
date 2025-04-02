package handlers

import (
	"database/sql"
	"encoding/json"
	"main/pkg/app/models"
	"net/http"
	"strconv"
)

func (h *Handler) ShowPost(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.PathValue("postId"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	userID := r.Context().Value("userID").(int)

	post, err := models.GetPost(h.db, postID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
