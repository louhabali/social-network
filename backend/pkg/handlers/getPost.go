package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *Handler) ShowPost(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(r.PathValue("postId"))
	post := GetPost(h.db, postID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
