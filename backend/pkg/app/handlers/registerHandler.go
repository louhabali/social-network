package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"main/pkg/app/models"
	"main/pkg/app/utils"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to decode the request body", http.StatusBadRequest)
		return
	}

	user := models.User{
		Firstname:   r.FormValue("firstName"),
		Lastname:    r.FormValue("lastName"),
		Email:       r.FormValue("email"),
		Password:    r.FormValue("password"),
		DateOfBirth: r.FormValue("dateOfBirth"),
		AboutMe:     r.FormValue("aboutMe"),
	}

	user.Nickname = utils.GetOptionalString(r.FormValue("nickname"))

	// Validate input
	if user.DateOfBirth == "" || user.Firstname == "" || user.Lastname == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Please fill all the fields", http.StatusBadRequest)
		return
	}

	if !utils.IsValidDateOfBirth(user.DateOfBirth) {
		http.Error(w, "You must be at least 5 years old", http.StatusBadRequest)
		return
	}

	if !utils.IsValidEmail(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if !utils.IsValidPassword(user.Password) {
		http.Error(w, "Password must be at least 8 characters and include at least one uppercase letter, one lowercase letter, and one number", http.StatusBadRequest)
		return
	}

	if !utils.IsValidName(user.Firstname) || !utils.IsValidName(user.Lastname) {
		http.Error(w, "Firstname and Lastname must contain only letters and be at most 10 characters long", http.StatusBadRequest)
		return
	}

	// Handle avatar upload
	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		extension := path.Ext(header.Filename)

		if header.Size > 1024*1024 {
			http.Error(w, "File too large. Max size is 1MB", http.StatusBadRequest)
			return
		}

		if !utils.IsValidAvatarExtension(extension) {
			http.Error(w, "Invalid avatar type. Allowed: JPEG, PNG", http.StatusBadRequest)
			return
		}

		rand, _ := utils.GenerateToken()
		name := rand + extension
		path := "uploads/avatars/" + name
		out, err := os.Create(path)
		if err != nil {
			http.Error(w, "Failed to save avatar", http.StatusInternalServerError)
			return
		}
		defer out.Close()
		io.Copy(out, file)
		user.Avatar = &path
	}

	// Insert user using models
	userID, err := user.Insert(h.db)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			http.Error(w, "email already used!", http.StatusBadRequest)
			return
		}
		if err.Error() == "UNIQUE constraint failed: users.first_name, users.last_name" {
			http.Error(w, "user already used!", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to insert the user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Success, User ID: %d", userID)))
}
