package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"main/pkg/app/models"
)

func (h *Handler) ShowPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := r.Context().Value("userID").(int)
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	groupId := r.URL.Query().Get("groupId")
	userId := r.URL.Query().Get("userId")

	var baseWhere string
	var args []interface{}

	if groupId != "" {
		baseWhere = " WHERE posts.group_id = ?"
		args = append(args, groupId)
	} else {
		baseWhere = " WHERE posts.group_id = 0"
		if userId != "" {
			baseWhere += " AND posts.user_id = ?"
			args = append(args, userId)
		}
	}

	posts, err := models.GetPostsWithPrivacy(userID, baseWhere, args, offset, h.db)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
