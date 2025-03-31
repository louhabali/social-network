package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"main/pkg/app/models"
	"main/pkg/app/utils"
)

// Login struct for request body
type Login struct {
	Email string `json:"email"`
	Pass  string `json:"password"`
}

// Handle user login
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var log Login
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if log.Email == "" || log.Pass == "" {
		http.Error(w, "Please fill all the fields", http.StatusBadRequest)
		return
	}

	// Fetch user from database
	user, err := models.GetUserByEmail(h.db, log.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	// Check password
	if !user.ComparePassword(log.Pass) {
		http.Error(w, "Incorrect password", http.StatusBadRequest)
		return
	}

	// Generate token
	uuid, err := utils.GenerateToken()
	if err != nil {
		http.Error(w, "Failed to generate session token", http.StatusInternalServerError)
		return
	}

	// Store token in database
	if err := models.StoreToken(h.db, user.ID, uuid); err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set HTTP cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Token",
		Value:    uuid,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(200 * 365 * 24 * time.Hour), // Long-lived token
	})

	w.Write([]byte("Connected successfully"))
}
