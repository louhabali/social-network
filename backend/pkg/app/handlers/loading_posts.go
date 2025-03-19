package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"
)

type Post struct {
	Pid             int
	Uid             int
	Username        string
	Avatar          *string
	Image           *string
	Content         string
	Privacy         string
	CreatedAt       string
	Likes           int
	Dislikes        int
	NbComment       int
	UserInteraction int
}

func (h *Handler) ShowPosts(w http.ResponseWriter, r *http.Request) {

	offset := r.URL.Query().Get("offset")
	groupId := r.URL.Query().Get("groupId")
	userId := r.URL.Query().Get("userId")

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		fmt.Println("invalid userId")
		return
	}

	query := `
	SELECT posts.id, posts.user_id, posts.created_at, posts.content, posts.image_path, posts.privacy, users.first_name, users.last_name, users.avatar
	FROM posts
	JOIN users ON posts.user_id = users.id
	`
	var args []interface{}
	if groupId != "" {
		query += " WHERE group_id = ?"
		args = append(args, groupId)
	} else {
		query += " WHERE group_id IS 0"
		if userId != "" {
			query += " AND user_id = ?"
			args = append(args, userId)
		}
	}
	fmt.Println(offset, groupId, userId)
	query += " ORDER BY posts.id DESC LIMIT 4 OFFSET ?"
	args = append(args, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var timeCreated time.Time
		var firstname, lastname string
		rows.Scan(&post.Pid, &post.Uid, &timeCreated, &post.Content, &post.Image, &post.Privacy, &firstname, &lastname, &post.Avatar)
		post.Content = html.UnescapeString(post.Content)
		post.Username = strings.Join([]string{firstname, lastname}, " ")
		post.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")
		if post.Avatar != nil {
			*post.Avatar, _ = EncodeImg(*post.Avatar)
		}
		h.db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", post.Pid).Scan(&post.NbComment)
		post.Likes = getPostLikeDisLike(post.Pid, 1, h.db)
		post.UserInteraction = getPostInteractions(post.Pid, userID, h.db)

		if post.Privacy == "almost_private" {
			isApproved := ""
			err = h.db.QueryRow("SELECT status FROM follow_requests WHERE (follower_id = ? AND following_id = ?) OR (following_id = ? AND follower_id = ?); ", post.Uid, userID, post.Uid, userID).Scan(&isApproved)

			if err == nil && isApproved == "approved" || post.Uid == userID {
				posts = append(posts, post)
			}
		} else if post.Privacy == "private" {
			exist := 0
			h.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM post_privacy_users WHERE post_id = ? AND user_id = ?)`, post.Pid, userID).Scan(&exist)
			if exist == 1 || post.Uid == userID {
				posts = append(posts, post)
			}

		} else if post.Privacy == "public" {
			posts = append(posts, post)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func getPostInteractions(pid, uid int, db *sql.DB) int {
	var ii int
	db.QueryRow(`select interaction from post_interactions where user_id = ? and post_id=?`, uid, pid).Scan(&ii)
	return ii
}

func GetPost(db *sql.DB, postID, userID int) Post {
	var p Post
	var timeCreated time.Time
	var fistname, lastname string
	db.QueryRow(`SELECT posts.id, posts.user_id, posts.created_at, posts.content,posts.image_path,users.first_name, users.last_name, users.avatar
        FROM posts
        JOIN users ON posts.user_id = users.id WHERE posts.id = ? `, postID).Scan(&p.Pid, &p.Uid, &timeCreated, &p.Content, &p.Image, &fistname, &lastname, &p.Avatar)
	p.Content = html.UnescapeString(p.Content)
	p.Username = strings.Join([]string{fistname, lastname}, " ")
	db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", p.Pid).Scan(&p.NbComment)
	p.Likes = getPostLikeDisLike(p.Pid, 1, db)
	p.UserInteraction = getPostInteractions(p.Pid, userID, db)
	if p.Avatar != nil {
		*p.Avatar, _ = EncodeImg(*p.Avatar)
	}
	p.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")
	return p
}

func getPostLikeDisLike(post_id, inter int, db *sql.DB) int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM post_interactions WHERE post_id=? AND interaction=?", post_id, inter).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}
