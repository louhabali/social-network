package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) ShowPost(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(r.PathValue("postId"))
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		fmt.Println("invalid userId")
		return
	}
	post := GetPost(h.db, postID,userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
