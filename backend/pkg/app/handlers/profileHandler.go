package handlers

import (
	"encoding/json"
	"main/pkg/app/models"
	"net/http"
	"strconv"
)

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	userID, _ := strconv.Atoi(r.PathValue("userId"))
	fid := r.Context().Value("userID").(int)

	userInfo, err := models.GetUserInfo(h.db, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	followStatus, err := models.GetFollowStatus(h.db, fid, userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	followers, err := models.GetFollowers(h.db, userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	whereClause := " WHERE user_id = ? AND group_id = 0"
	posts, err := models.GetPostsWithPrivacy(fid, whereClause, []interface{}{userID}, 0, h.db)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	following, err := models.GetFollowing(h.db, userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	totalPosts, err := models.GetUserPostCount(h.db, userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := struct {
		User         models.User       `json:"user"`
		Posts        []models.Post     `json:"posts"`
		Followers    []models.Follower `json:"followers"`
		Following    []models.Follower `json:"following"`
		UserId       int               `json:"userid"`
		FollowStatus string            `json:"followStatus"`
		TotalPosts   int               `json:"totalPosts"`
	}{
		User:         userInfo,
		Posts:        posts,
		Followers:    followers,
		Following:    following,
		UserId:       fid,
		FollowStatus: followStatus,
		TotalPosts:   totalPosts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
