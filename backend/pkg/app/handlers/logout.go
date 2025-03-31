package handlers

import (
	"net/http"
	"main/pkg/app/models"
)

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Token")
	if err != nil {
		http.Error(w, "No token found", http.StatusBadRequest)
		return
	}

	// Delete cookie1
	http.SetCookie(w, &http.Cookie{
		Name:   "Token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	err = models.RemoveToken(h.db, cookie.Value)
	if err != nil {
		http.Error(w, "Failed to update the token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
