package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int
	Nickname    *string
	Firstname   string
	Lastname    string
	Email       string
	Password    string
	DateOfBirth string
	AboutMe     string
	Avatar      *string
	IsPublic    bool
}

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// decode the request body
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("failed to decode the request body"))

			fmt.Println("failed to decode the request body", err)
			return
		}

		reg := User{
			Firstname:   r.FormValue("firstName"),
			Lastname:    r.FormValue("lastName"),
			Email:       r.FormValue("email"),
			Password:    r.FormValue("password"),
			DateOfBirth: r.FormValue("dateOfBirth"),
			AboutMe:     r.FormValue("aboutMe"),
		}

		// check if the informations are not empty
		if reg.DateOfBirth == "" || reg.Firstname == "" || reg.Lastname == "" || reg.Email == "" || reg.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Please fill all the fields"))
			return
		}

		if nickname := r.FormValue("nickname"); nickname != "" {
			reg.Nickname = &nickname
		}
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(reg.Email) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid email format"))
			return
		}
		
		if !isValidPassword(reg.Password) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Password must be at least 8 characters ,include at least one uppercase letter, one lowercase letter, and one number"))
			return
		}
		nameRegex := regexp.MustCompile(`^[a-zA-Z]{1,10}$`)
		if !nameRegex.MatchString(reg.Firstname) || !nameRegex.MatchString(reg.Lastname) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Firstname and Lastname must contain only letters and be at most 10 characters long"))
			return
		}
		dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
		if !dateRegex.MatchString(reg.DateOfBirth) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Date of Birth must be in DD-MM-YYYY format"))
			return
		}

		// hash the password first
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("failed to hash the password", err)
			return
		}

		file, header, err := r.FormFile("avatar")
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
			avatarPath := "uploads/avatars/" + name
			out, err := os.Create(avatarPath)
			if err != nil {

				http.Error(w, "Failed to save avatar", http.StatusInternalServerError)
				fmt.Println("66", err)
				return
			}
			defer out.Close()
			io.Copy(out, file)
			reg.Avatar = &avatarPath
		}
		fmt.Println(reg.Avatar)
		// insert them to users table these informations
		var id int
		err = h.db.QueryRow(`
		INSERT INTO users 
		(nickname,date_of_birth,about_me,first_name,last_name,email,password_hash,avatar) 
		VALUES 
		(?,?,?,?,?,?,?,?) RETURNING id`, reg.Nickname, reg.DateOfBirth, reg.AboutMe, reg.Firstname, reg.Lastname, reg.Email, hashedPass, reg.Avatar).Scan(&id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("failed to insert the user"))

			fmt.Println("failed to insert the user", err)
			return
		}

		fmt.Println("000")
		w.Write([]byte("success"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func GenereteTocken() (string, error) {
	u2, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
		return "", err
	}
	return u2.String(), nil
}
func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	return hasUpper && hasLower && hasDigit
}