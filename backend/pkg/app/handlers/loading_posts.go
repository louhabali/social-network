package handlers

import (
	"database/sql"
	"encoding/json"
	"html"
	"net/http"
	"strconv"
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

func (h *Handler) getPostsWithPrivacy(userID int, baseWhere string, args []interface{},offset int) ([]Post, error) {
	privacyWhere := `
        AND (
            posts.privacy = 'public' 
            OR posts.user_id = ? 
            OR (posts.privacy = 'almost_private' AND EXISTS (
                SELECT 1 FROM follow_requests 
                WHERE follower_id = ? AND following_id = posts.user_id AND status = 'approved'
            )) 
            OR (posts.privacy = 'private' AND EXISTS (
                SELECT 1 FROM post_privacy_users 
                WHERE post_id = posts.id AND user_id = ?
            ))
        )`

	args = append(args, userID, userID, userID)

	query := `
        SELECT posts.id, posts.user_id, posts.created_at, posts.content, 
               posts.image_path, posts.privacy, 
               users.first_name, users.last_name, users.avatar,
               (SELECT COUNT(*) FROM comments WHERE post_id = posts.id) AS comment_count
        FROM posts
        JOIN users ON posts.user_id = users.id
        ` + baseWhere + privacyWhere + `
        ORDER BY posts.id DESC LIMIT 4 OFFSET ?`
	args = append(args, offset)
	rows, err := h.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var timeCreated time.Time
		var firstName, lastName string

		err := rows.Scan(&post.Pid, &post.Uid, &timeCreated, &post.Content,
			&post.Image, &post.Privacy, &firstName, &lastName, &post.Avatar, &post.NbComment)
		if err != nil {
			return nil, err
		}

		post.Content = html.UnescapeString(post.Content)
		post.Username = firstName + " " + lastName
		post.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")

		if post.Avatar != nil {
			*post.Avatar, _ = EncodeImg(*post.Avatar)
		}

		post.Likes = getPostLikeDisLike(post.Pid, 1, h.db)
		post.UserInteraction = getPostInteractions(post.Pid, userID, h.db)

		posts = append(posts, post)
	}

	return posts, nil
}

func (h *Handler) ShowPosts(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	posts, err := h.getPostsWithPrivacy(userID, baseWhere, args,offset)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
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
