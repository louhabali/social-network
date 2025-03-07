package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *Handler) ShowPost(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(int)
	postID, _ := strconv.Atoi(r.PathValue("postId"))
	post := GetPost(h.db, postID)
	post.Likes = getPostLikeDisLike(postID, 1, h.db)
	post.UserInteraction = getPostInteractions(postID, userID, h.db)
	data := struct {
		Post            Post
		Likes int
		UserInteraction int
	}{
		Post:            post,
		Likes: post.Likes,
		UserInteraction: post.UserInteraction,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
