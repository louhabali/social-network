package handlers

import (
	"encoding/json"
	"io"
	"main/pkg/app/models"
	"main/pkg/app/utils"
	"net/http"
	"os"
	"path"
	"strconv"
)

func (h *Handler) CommentCreationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		postID, _ := strconv.Atoi(r.PathValue("postId"))
		userID := r.Context().Value("userID").(int)
		content := r.FormValue("content")

		// Handle image upload
		file, header, err := r.FormFile("image")
		if content == "" && file == nil {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		var imagePath *string
		if err == nil { // If an image is uploaded
			defer file.Close()
			fileExt := path.Ext(header.Filename)

			// Validate image format and size
			if err := utils.ValidatePostImage(fileExt, header.Size); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Generate unique file name
			rand, _ := utils.GenerateToken()
			imageName := rand + fileExt
			fullPath := "uploads/images/" + imageName

			// Save image to disk
			out, err := os.Create(fullPath)
			if err != nil {
				http.Error(w, "Failed to save image", http.StatusInternalServerError)
				return
			}
			defer out.Close()
			io.Copy(out, file)

			imagePath = &fullPath
		}

		// Insert comment
		commentID, err := models.InsertComment(h.db, postID, userID, content, imagePath)
		if err != nil {
			http.Error(w, "Failed to insert the comment", http.StatusInternalServerError)
			return
		}

		// Fetch and return the new comment
		comment, err := models.GetComment(h.db, commentID)
		if err != nil {
			http.Error(w, "Failed to fetch the comment", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comment)

	} else if r.Method == http.MethodGet {
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		postID, _ := strconv.Atoi(r.PathValue("postId"))

		comments, err := models.GetCommentsByPost(h.db, postID, offset)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
