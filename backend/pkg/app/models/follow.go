package models

import (
	"database/sql"
	"fmt"
)

type Follower struct {
	ID        int     `json:"id"`
	Firstname string  `json:"fname"`
	Lastname  string  `json:"lname"`
	Avatar    *string `json:"avatar"`
}

func GetFollowers(db *sql.DB, userID int) ([]Follower, error) {
	rows, err := db.Query(`
		SELECT users.id, users.first_name, users.last_name FROM users
		JOIN follow_requests ON users.id = follow_requests.follower_id
		WHERE follow_requests.following_id = ? AND follow_requests.status = 'approved';
	`, userID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var followers []Follower
	for rows.Next() {
		var f Follower
		err := rows.Scan(&f.ID, &f.Firstname, &f.Lastname)
		if err != nil {
			return nil, err
		}
		followers = append(followers, f)
	}
	return followers, nil
}

func GetFollowing(db *sql.DB, userID int) ([]Follower, error) {
	rows, err := db.Query(`
		SELECT users.id, users.first_name, users.last_name FROM users
		JOIN follow_requests ON users.id = follow_requests.following_id
		WHERE follow_requests.follower_id = ? AND follow_requests.status = 'approved';
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []Follower
	for rows.Next() {
		var f Follower
		err := rows.Scan(&f.ID, &f.Firstname, &f.Lastname)
		if err != nil {
			return nil, err
		}
		following = append(following, f)
	}
	return following, nil
}

func GetFollowStatus(db *sql.DB, followerID, followingID int) (string, error) {
	var status string
	err := db.QueryRow("SELECT status FROM follow_requests WHERE follower_id = ? AND following_id = ?", followerID, followingID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return status, nil
}


func NotFollowing(db *sql.DB, userID int) ([]Follower,error) {
	rows, err := db.Query(`SELECT users.id, users.first_name, users.last_name ,users.avatar FROM users
    LEFT JOIN follow_requests ON users.id = follow_requests.follower_id
    WHERE follow_requests.follower_id != ? or users.id !=?`, userID, userID)
	if err != nil {
		return nil , err
	}
	defer rows.Close()

	var notFollowing []Follower
	for rows.Next() {
		var nf Follower
		err := rows.Scan(&nf.ID, &nf.Firstname, &nf.Lastname, &nf.Avatar)
		if err != nil {
			return nil, err
		}
		notFollowing = append(notFollowing, nf)
	}
	return notFollowing , nil
}

// InsertFollowRequest inserts a new follow request (Pending or Approved)
func InsertFollowRequest(db *sql.DB, followerID, followingID int, status string) (int64, error) {

	_,err := GetUserInfo(db,followingID)
	if err != nil {
		return 0 , err
	}

	res, err := db.Exec("INSERT INTO follow_requests (follower_id, following_id, status) VALUES (?, ?, ?)", followerID, followingID, status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// DeleteFollowRequest removes a follow request
func DeleteFollowRequest(db *sql.DB, followerID, followingID int) error {
	_, err := db.Exec("DELETE FROM follow_requests WHERE follower_id = ? AND following_id = ?", followerID, followingID)
	return err
}

// ApproveFollowRequest updates a follow request's status to 'approved'
func ApproveFollowRequest(db *sql.DB, requestID int) error {
	_, err := db.Exec("UPDATE follow_requests SET status = 'approved' WHERE id = ?", requestID)
	return err
}
