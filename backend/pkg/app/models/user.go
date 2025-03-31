package models

import (
	"database/sql"
	"errors"

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

// Insert a new user into the database
func (u *User) Insert(db *sql.DB) (int, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.New("failed to hash password")
	}

	var id int
	err = db.QueryRow(`
		INSERT INTO users 
		(nickname, date_of_birth, about_me, first_name, last_name, email, password_hash, avatar) 
		VALUES (?,?,?,?,?,?,?,?) RETURNING id`,
		u.Nickname, u.DateOfBirth, u.AboutMe, u.Firstname, u.Lastname, u.Email, hashedPass, u.Avatar,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}


func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	user := &User{}
	err := db.QueryRow(`SELECT id, password_hash FROM users WHERE email = ?`, email).
		Scan(&user.ID, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func GetUserInfo(db *sql.DB, userID int) (User, error) {
	var user User
	err := db.QueryRow("select id,email,nickname,first_name,last_name,date_of_birth,avatar,about_me,is_public from users where id = ?", userID).
		Scan(&user.ID, &user.Email, &user.Nickname, &user.Firstname, &user.Lastname, &user.DateOfBirth, &user.Avatar, &user.AboutMe, &user.IsPublic)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// CheckIfPublic checks if a user has a public profile
func CheckIfPublic(db *sql.DB, userID int) (bool, error) {
	var isPublic bool
	err := db.QueryRow("SELECT is_public FROM users WHERE id = ?", userID).Scan(&isPublic)
	if err != nil {
		return false, err
	}
	return isPublic, nil
}


// ToggleUserPrivacy updates the user's privacy setting
func ToggleUserPrivacy(db *sql.DB, userID int, isPublic bool) error {
	_, err := db.Exec("UPDATE users SET is_public = ? WHERE id = ?", isPublic, userID)
	return err
}


func GetUserIDFromToken(uuid string, db *sql.DB) (int, error) {
	var id int
	err := db.QueryRow("SELECT user_id FROM tokens WHERE token=?", uuid).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}