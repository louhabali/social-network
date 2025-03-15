package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"

	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Email string `json:"email"`
	Pass  string `json:"password"`
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var log Login

		json.NewDecoder(r.Body).Decode(&log)
		if log.Email == "" || log.Pass == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("please fill all the fields"))
			return
		}
		var passScanned string
		var UserID int
		err := h.db.QueryRow(`select password_hash,id from users where email = ?`, log.Email).Scan(&passScanned, &UserID)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("user not found"))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("database error"))
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(passScanned), []byte(log.Pass))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("incorrect password"))

			fmt.Println("failed to compare the db password with db password")
			return
		}
		uuid, err := uuid.NewV4()
		if err != nil {
			fmt.Println("failed to generate a uuid")
			return
		}

		_, err = h.db.Exec(`INSERT INTO tokens(user_id,token) VALUES (?,?)`, UserID,uuid.String())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to insert the user"))

			fmt.Println("failed to insert the user", err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "Token",
			Value:    uuid.String(),
			HttpOnly: true,
			Path: "/",
<<<<<<< HEAD
			Expires:  time.Now().Add(365 * 24 * time.Hour),
=======
			Expires:  time.Now().Add(200 * 365 * 24 * time.Hour),
>>>>>>> alouhab
		})

		w.Write([]byte("connected successfully"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
