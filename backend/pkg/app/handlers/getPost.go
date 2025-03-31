package handlers

import (
	"encoding/json"
	"main/pkg/app/models"
	"net/http"
	"strconv"
)

func (h *Handler) ShowPost(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(r.PathValue("postId"))
	userID := r.Context().Value("userID").(int)

	post := models.GetPost(h.db, postID, userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
