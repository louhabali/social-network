package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"main/pkg/app/models"
	"main/pkg/app/utils"
	"net/http"
	"os"
	"path"
	"strconv"
)

// Handle post creation
func (h *Handler) PostCreationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)

	content := r.FormValue("content")
	privacy := r.FormValue("privacy")
	groupID, _ := strconv.Atoi(r.FormValue("groupId"))
	
	file, header, err := r.FormFile("image")

	if content == "" && file == nil {
		http.Error(w, "Content or image is required", http.StatusBadRequest)
		return
	}

	if len(content) > 200 {
		http.Error(w, "Content too long", http.StatusBadRequest)
		return
	}

	if privacy != "public" && privacy != "private" && privacy != "almost_private" {
		http.Error(w, "Invalid privacy setting", http.StatusBadRequest)
		return
	}

	if groupID != 0 && privacy != "public" {
		http.Error(w, "Group posts must be public", http.StatusBadRequest)
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

	// Create post in database
	postID, err := models.CreatePost(h.db, userID, content, privacy, imagePath, groupID)
	if err == sql.ErrNoRows {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Failed to insert post", http.StatusInternalServerError)
		return
	}

	// Process selected users for privacy settings
	var selectedUsers []int
	for _, sID := range r.Form["selectedUsers"] {
		i, _ := strconv.Atoi(sID)
		selectedUsers = append(selectedUsers, i)
	}
	if err := models.AddPostPrivacy(h.db, postID, selectedUsers); err != nil {
		http.Error(w, "Failed to set post privacy", http.StatusInternalServerError)
		return
	}
	p, err := models.GetPost(h.db, postID, userID)
	if err != nil {
		http.Error(w, "Failed to retrieve post", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
