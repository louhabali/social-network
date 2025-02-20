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
	Avatar          string
	Image           string
	Content         string
	CreatedAt       string
	Likes           int
	Dislikes        int
	NbComment       int
	UserInteraction int
}

func (h *Handler) ShowPosts(w http.ResponseWriter, r *http.Request) {

	var data struct {
		Offset   int    `json:"offset"`
	}
	json.NewDecoder(r.Body).Decode(&data)

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		fmt.Println("invalid userId")
		return
	}
	//The user must be able to specify the privacy of the post:

	//public (all users in the social network will be able to see the post)
	//almost private (only followers of the creator of the post will be able to see the post)
	//private (only the followers chosen by the creator of the post will be able to see it)
	//The user must be able to specify the type of post:
	query := `
        SELECT posts.id, posts.user_id, posts.created_at, posts.content,posts.privacy,users.fistname, users.lastname 
        FROM posts
        JOIN users ON posts.user_id = users.id`

	rows, _ := h.db.Query(query)
	var posts []Post
	for rows.Next() {
		var post Post
		var timeCreated time.Time
		var fistname, lastname ,privacy string
		rows.Scan(&post.Pid, &post.Uid,  &timeCreated, &post.Content,&privacy, &fistname, &lastname)
		post.Username = strings.Join([]string{fistname, lastname}, " ")
		post.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")

		h.db.QueryRow("SELECT COUNT(*) FROM Comments WHERE PostID = ?", post.Pid).Scan(&post.NbComment)
		post.Likes = getPostLikeDisLike(post.Pid, 1, h.db)
		post.Dislikes = getPostLikeDisLike(post.Pid, -1, h.db)
		post.UserInteraction = getPostInteractions(post.Pid, userID, h.db)
		if privacy == "public" {
			posts = append(posts, post)
		} else if privacy == "only followers" {
			var isfollower int
			err := h.db.QueryRow("select COUNT(*) from follow_requests where follower_id = ? and following_id = ? and status='accepted' ",userID, post.Pid).Scan(&isfollower)
			if err != nil {
				fmt.Println("not a follower")
				continue
			} else {
				posts = append(posts, post)
			}
		} else {
			// hadi ta nchoufo smo ndiro wash checkbox wla mn3rf
		}


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
	db.QueryRow(`SELECT posts.id, posts.user_id, posts.created_at, posts.content,users.firstname, users.lastname 
        FROM posts
        JOIN users ON posts.user_id = users.id WHERE posts.id = ? `, postID).Scan(&p.Pid, &p.Uid, &timeCreated, &p.Content, &fistname, &lastname)
	p.Username = strings.Join([]string{fistname, lastname}, " ")
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
