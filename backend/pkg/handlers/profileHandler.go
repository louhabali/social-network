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
	"time"
)

type Follower struct {
	ID        int
	Firstname string
	Lastname  string
	Nickname  string
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.PathValue("userId"))
	fid, ok := r.Context().Value("userID").(int)
	
	if !ok {
		fmt.Println("invalid userId")
		return
	}

	posts, err := getPostsByUserID(h.db, userID)
	if err != nil {
		fmt.Println("Error fetching posts:", err)
	}

	followers, err := getFollowers(h.db, userID)
	if err != nil {
		fmt.Println("Error fetching followers:", err)
	}

	following, err := getFollowing(h.db, userID)
	if err != nil {
		fmt.Println("Error fetching following:", err)
	}

	userInfo, err := getUserInfo(h.db, userID)
	if err != nil {
		fmt.Println("Error fetching profile privacy:", err)
	}
	c := 0
	err = h.db.QueryRow("SELECT COUNT(*) FROM users JOIN follow_requests ON users.id = follow_requests.following_id WHERE follow_requests.follower_id = ? AND follow_requests.following_id = ?", userID, fid).Scan(&c)
	isFollowing := c != 0

	response := struct {
		User        User       `json:"user"`
		Posts       []Post     `json:"posts"`
		Followers   []Follower `json:"followers"`
		Following   []Follower `json:"following"`
		IsFollowing bool       `json:"isFollowing"`
		UserId      int        `json:"userid"`
	}{
		User:        userInfo,
		Posts:       posts,
		Followers:   followers,
		Following:   following,
		IsFollowing: isFollowing,
		UserId:      fid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getPostsByUserID(db *sql.DB, userID int) ([]Post, error) {
	rows, _ := db.Query(`SELECT posts.id, posts.user_id, posts.created_at, posts.content,posts.image_path,users.first_name, users.last_name, users.avatar
        FROM posts
        JOIN users ON posts.user_id = users.id  WHERE posts.user_id = ? ORDER BY posts.id DESC LIMIT ? OFFSET ?`, userID, 4, 0)

	var posts []Post
	for rows.Next() {
		var post Post
		var timeCreated time.Time
		var firstName, lastName string
		rows.Scan(&post.Pid, &post.Uid, &timeCreated, &post.Content, &post.Image, &firstName, &lastName, &post.Avatar)
		post.Username = firstName + " " + lastName
		post.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")
		if post.Avatar != nil {
			*post.Avatar, _ = EncodeImg(*post.Avatar)
		}
		db.QueryRow("SELECT COUNT(*) FROM Comments WHERE PostID = ?", post.Pid).Scan(&post.NbComment)
		post.Likes = getPostLikes(post.Pid, 1, db)
		post.UserInteraction = getPostInteractions(post.Pid, userID, db)
		posts = append(posts, post)
	}
	return posts, nil
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
	`, userID)

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
