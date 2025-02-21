package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
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
		userId, err = GetUserIDFromToken(c.Value, db)
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

func GetUserIDFromToken(uuid string, db *sql.DB) (int, error) {
	var id int
	err := db.QueryRow("SELECT user_id FROM tokens WHERE token=?", uuid).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
