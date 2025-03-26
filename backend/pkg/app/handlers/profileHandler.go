package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Follower struct {
	ID        int     `json:"id"`
	Firstname string  `json:"fname"`
	Lastname  string  `json:"lname"`
	Nickname  string  `json:"nickname"`
	Avatar    *string `json:"avatar"`
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.PathValue("userId"))
	fid, ok := r.Context().Value("userID").(int)
	if !ok {
		fmt.Println("invalid userId")
		return
	}

	userInfo, err := getUserInfo(h.db, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		fmt.Println("Error fetching profile privacy:", err)
		return
	}

	followStatus := ""
	err = h.db.QueryRow("SELECT status FROM follow_requests WHERE follower_id = ? AND following_id = ?", fid, userID).Scan(&followStatus)
	fmt.Println(followStatus, fid, userID, err)
	followers, err := getFollowers(h.db, userID)
	if err != nil {
		fmt.Println("Error fetching followers:", err)
	}

	whereClause := " WHERE user_id = ? AND group_id = 0"
	posts, err := h.getPostsWithPrivacy(fid, whereClause, []interface{}{userID},0)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	following, err := getFollowing(h.db, userID)
	if err != nil {
		fmt.Println("Error fetching following:", err)
	}

	totalPosts, err := getUserPostCount(h.db,userID)
	if err != nil {
		fmt.Println("Error fetching total posts:", err)
	}
	response := struct {
		User         User       `json:"user"`
		Posts        []Post     `json:"posts"`
		Followers    []Follower `json:"followers"`
		Following    []Follower `json:"following"`
		UserId       int        `json:"userid"`
		FollowStatus string     `json:"followStatus"`
		TotalPosts   int        `json:"totalPosts"`
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

func getUserPostCount(db *sql.DB, userID int) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ? AND group_id = 0", userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func getUserInfo(db *sql.DB, userID int) (User, error) {
	var user User
	err := db.QueryRow("select id,email,nickname,first_name,last_name,date_of_birth,avatar,about_me,is_public from users where id = ?", userID).
		Scan(&user.ID, &user.Email, &user.Nickname, &user.Firstname, &user.Lastname, &user.DateOfBirth, &user.Avatar, &user.AboutMe, &user.IsPublic)
	if err != nil {
		return User{}, err
	}
	if user.Avatar != nil {
		*user.Avatar, _ = EncodeImg(*user.Avatar)
	}
	return user, nil
}

func getFollowers(db *sql.DB, userID int) ([]Follower, error) {
	rows, err := db.Query(`
		SELECT users.id, users.first_name, users.last_name, users.nickname FROM users
		JOIN follow_requests ON users.id = follow_requests.follower_id
		WHERE follow_requests.following_id = ? AND follow_requests.status = 'approved';
	`, userID)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	var followers []Follower
	for rows.Next() {
		var f Follower
		rows.Scan(&f.ID, &f.Firstname, &f.Lastname, &f.Nickname)
		followers = append(followers, f)
	}
	return followers, nil
}

func getFollowing(db *sql.DB, userID int) ([]Follower, error) {
	rows, _ := db.Query(`
		SELECT users.id, users.first_name, users.last_name, users.nickname FROM users
		JOIN follow_requests ON users.id = follow_requests.following_id
		WHERE follow_requests.follower_id = ? AND follow_requests.status = 'approved';
	`, userID, userID)

	var following []Follower
	for rows.Next() {
		var f Follower
		rows.Scan(&f.ID, &f.Firstname, &f.Lastname, &f.Nickname)
		following = append(following, f)
	}
	return following, nil
}

func EncodeImg(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}

	if len(imageBytes) == 0 {
		return "", fmt.Errorf("image file is empty")
	}

	encodedImage := base64.StdEncoding.EncodeToString(imageBytes)
	return encodedImage, nil
}
