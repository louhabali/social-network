package models

import "database/sql"

// InvalidateToken removes the token from the database
func RemoveToken(db *sql.DB, token string) error {
	_, err := db.Exec(`DELETE FROM tokens WHERE token=?`, token)
	return err
}


func StoreToken(db *sql.DB, userID int, token string) error {
	_, err := db.Exec(`INSERT INTO tokens(user_id, token) VALUES (?, ?)`, userID, token)
	return err
}