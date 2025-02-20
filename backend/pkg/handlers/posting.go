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
)

func PostCreationHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			fmt.Println("invalid userId")
			return
		}

		p := Post{
			Content: r.FormValue("content"),
		}

		if p.Content == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("all fields are required")
			return
		}

		file, header, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			rand, _ := GenereteTocken()
			extensions := path.Ext(header.Filename)
			name := rand + extensions
			imagePath := "uploads/" + name
			out, err := os.Create(imagePath)
			if err != nil {
				http.Error(w, "Failed to save image", http.StatusInternalServerError)
				return
			}
			defer out.Close()
			io.Copy(out, file)
			p.Image = imagePath
		}

		postID := 0

		err = db.QueryRow(`INSERT INTO posts (user_id, content, image_path) VALUES (?, ?, ?) RETURNING id`, userID, html.EscapeString(p.Content), p.Image).Scan(&postID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("failed to insert the post")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GetPost(db, postID))

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
