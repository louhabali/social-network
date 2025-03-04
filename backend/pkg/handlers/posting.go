package handlers

import (
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

func (h *Handler) PostCreationHandler(w http.ResponseWriter, r *http.Request) {
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
			Privacy: r.FormValue("privacy"),
		}

		groupID, _ := strconv.Atoi(r.FormValue("groupId"))

		file, header, err := r.FormFile("image")

		if p.Content == "" && file == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("all fields are required")
			return
		}

		if err == nil {
			defer file.Close()
			extensions := path.Ext(header.Filename)
			if !IsValidImageExtension(extensions) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("Invalid image type. Supported formats: JPEG, PNG, GIF")
				return
			}
			rand, _ := GenereteTocken()
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

		postID := 0
		
		

	
		err = h.db.QueryRow(`INSERT INTO posts (user_id, content,privacy, image_path, group_id) VALUES (?, ?, ?, ?, ?) RETURNING id`, userID, html.EscapeString(p.Content),p.Privacy, p.Image, groupID).Scan(&postID)
		if err != nil {
			fmt.Println("err01 ",err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("failed to insert the post")
			return
		}
		selectUsers := r.Form["selectedUsers"]
		if len(selectUsers)>0 {
			for _, sID := range selectUsers {
				for _,ii := range strings.Split(sID,","){
					i , _:= strconv.Atoi(ii)
					_,err = h.db.Exec(`INSERT INTO post_privacy_users (user_id, post_id, created_at) VALUES (?, ?, ?)`, i, postID, time.Now())
					if err != nil {
						fmt.Println("err02 ",err)
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode("failed to insert post privacy")
						return
					}
				}
			 }
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GetPost(h.db, postID))

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func IsValidImageExtension(fileExt string) bool {
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	lowerFileExt := strings.ToLower(fileExt)
	for _, ext := range allowedExtensions {
		if ext == lowerFileExt {
			return true
		}
	}
	return false
}
