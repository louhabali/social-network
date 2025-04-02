package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"main/pkg/app/models"
	"net/http"
)

func (h *Handler) Midd(n http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uu, b := CheckIfCookieValid(w, r, h.db)
		if !b {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Logged out")
			return
		}
		ctx := context.WithValue(r.Context(), "userID", uu)
		n(w, r.WithContext(ctx))
	}
}

func CheckIfCookieValid(w http.ResponseWriter, r *http.Request, db *sql.DB) (int, bool) {
	var userId int
	c, err := r.Cookie("Token")
	if err == nil {
		userId, err = models.GetUserIDFromToken(c.Value, db)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:   "Token",
				MaxAge: -1,
			})
			return userId, false
		} else {
			return userId, true
		}
	}
	return userId, false
}

