package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Comment struct {
	Id        int
	Username  string
	Avatar    *string
	Image     *string
	Content   string
	CreatedAt string
}

func (h *Handler) CommentCreationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		postID, _ := strconv.Atoi(r.PathValue("postId"))
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			fmt.Println("invalid userId")
			return
		}
		p := Post{
			Content: r.FormValue("content"),
		}
		file, header, err := r.FormFile("image")

		if p.Content == "" && file == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("all fields are required")
			return
		}

		if err == nil {
			defer file.Close()
			rand, _ := GenereteTocken()
			extensions := path.Ext(header.Filename)
			name := rand + extensions
			imagePath := "uploads/images/" + name
			out, err := os.Create(imagePath)
			if err != nil {
				http.Error(w, "Failed to save image", http.StatusInternalServerError)
				return
			}
			defer out.Close()
			io.Copy(out, file)
			p.Image = &imagePath
		}

		commentId := 0

		err = h.db.QueryRow(`INSERT INTO comments (post_id,user_id, content, image_path) VALUES (?, ?, ?, ?) RETURNING id`, postID, userID, html.EscapeString(p.Content), p.Image).Scan(&commentId)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("failed to insert the post")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GetComment(h.db, commentId))

	} else if r.Method == http.MethodGet {
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		postID, _ := strconv.Atoi(r.PathValue("postId"))

		rows, err := h.db.Query(`SELECT comments.id, comments.content, comments.image_path, users.first_name, users.last_name, comments.created_at,users.avatar FROM comments
			JOIN users ON comments.user_id = users.id WHERE comments.post_id = ? ORDER BY comments.id DESC LIMIT 4 OFFSET ?;`, postID, offset)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var comments []Comment
		for rows.Next() {
			var comment Comment
			var timeCreated time.Time
			var fistname, lastname string
			rows.Scan(&comment.Id, &comment.Content, &comment.Image, &fistname, &lastname, &timeCreated, &comment.Avatar)
			comment.Content = html.UnescapeString(comment.Content)
			comment.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")
			comment.Username = strings.Join([]string{fistname, lastname}, " ")
			if comment.Avatar != nil {
				*comment.Avatar, _ = EncodeImg(*comment.Avatar)
			}
			comments = append(comments, comment)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func GetComment(db *sql.DB, commentID int) Comment {

	var c Comment
	var timeCreated time.Time
	var fistname, lastname string
	db.QueryRow(`SELECT comments.id ,comments.content, comments.image_path, users.first_name, users.last_name, comments.created_at,users.avatar FROM comments
			JOIN users ON comments.user_id = users.id WHERE comments.id = ? `, commentID).Scan(&c.Id, &c.Content, &c.Image, &fistname, &lastname, &timeCreated, &c.Avatar)
	c.Username = strings.Join([]string{fistname, lastname}, " ")
	c.Content = html.UnescapeString(c.Content)
	if c.Avatar != nil {
		*c.Avatar, _ = EncodeImg(*c.Avatar)
	}
	c.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")
	return c
}
