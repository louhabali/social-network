package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	CreatedAt       string
	Likes           int
	Dislikes        int
	NbComment       int
	UserInteraction int
}

func (h *Handler) ShowPosts(w http.ResponseWriter, r *http.Request) {

	var data struct {
		Offset int `json:"offset"`
	}

	json.NewDecoder(r.Body).Decode(&data)

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		fmt.Println("invalid userId")
		return
	}

	query := `
        SELECT posts.id, posts.user_id, posts.created_at, posts.content,posts.image_path,users.first_name, users.last_name , users.avatar
        FROM posts
        JOIN users ON posts.user_id = users.id ORDER BY posts.id DESC
        `

	//WHERE posts.privacy = 'public'
	//OR (posts.privacy = 'almost_private' AND posts.user_id IN (SELECT follower_id FROM followers WHERE user_id = ?))
	//OR (posts.privacy = 'private' AND posts.id IN (SELECT post_id FROM post_privacy WHERE user_id = ?))
	rows, err := h.db.Query(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	var posts []Post
	for rows.Next() {
		var post Post
		var timeCreated time.Time
		var fistname, lastname string
		rows.Scan(&post.Pid, &post.Uid, &timeCreated, &post.Content, &post.Image,&fistname, &lastname, &post.Avatar)
		post.Username = strings.Join([]string{fistname, lastname}, " ")
		post.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")
		if post.Avatar != nil {
			*post.Avatar,_  = EncodeImg(*post.Avatar)
		}
		h.db.QueryRow("SELECT COUNT(*) FROM Comments WHERE PostID = ?", post.Pid).Scan(&post.NbComment)
		post.Likes = getPostLikeDisLike(post.Pid, 1, h.db)
		post.Dislikes = getPostLikeDisLike(post.Pid, -1, h.db)
		post.UserInteraction = getPostInteractions(post.Pid, userID, h.db)
		posts = append(posts, post)

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)

}

func getPostInteractions(pid, uid int, db *sql.DB) int {
	var ii int
	db.QueryRow(`select Interaction from PostInteractions where UserID = ? and PostID=?`, uid, pid).Scan(&ii)
	return ii
}

func GetPost(db *sql.DB, postID int) Post {

	var p Post
	var timeCreated time.Time
	var fistname, lastname string
	db.QueryRow(`SELECT posts.id, posts.user_id, posts.created_at, posts.content,posts.image_path,users.first_name, users.last_name, users.avatar
        FROM posts
        JOIN users ON posts.user_id = users.id WHERE posts.id = ? `, postID).Scan(&p.Pid, &p.Uid, &timeCreated, &p.Content, &p.Image,&fistname, &lastname, &p.Avatar)
	p.Username = strings.Join([]string{fistname, lastname}, " ")
	if p.Avatar != nil {
		*p.Avatar,_  = EncodeImg(*p.Avatar)
	}
	p.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")
	return p
}

func getPostLikeDisLike(post_id, inter int, db *sql.DB) int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM PostInteractions WHERE PostID=? AND Interaction=?", post_id, inter).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}
