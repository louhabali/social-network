package models

import (
	"database/sql"
	"main/pkg/app/utils"
)

type GroupEvent struct {
	ID            int    `json:"id"`
	GroupID       int    `json:"group_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Date          string `json:"date"`
	UserID        int    `json:"created_id"`
	CreatedAt     string `json:"created_at"`
	GoingCount    int    `json:"going_count"`
	NotGoingCount int    `json:"not_going_count"`
	UserResponse  string `json:"user_response"`
}

// GetGroupEvents retrieves events for a group with pagination
func GetGroupEvents(db *sql.DB, groupID int, offset int, userID int) ([]GroupEvent, error) {
	rows, err := db.Query("SELECT id, group_id, title, description, event_date, created_id, created_at FROM group_events WHERE group_id = ? ORDER BY id DESC LIMIT 4 OFFSET ?", groupID, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []GroupEvent
	for rows.Next() {
		var e GroupEvent
		err := rows.Scan(&e.ID, &e.GroupID, &e.Title, &e.Description, &e.Date, &e.UserID, &e.CreatedAt)
		if err != nil {
			return nil, err
		}

		db.QueryRow("SELECT COUNT(*) FROM event_responses WHERE event_id = ? AND response = 'going'", e.ID).Scan(&e.GoingCount)
		db.QueryRow("SELECT COUNT(*) FROM event_responses WHERE event_id = ? AND response = 'not_going'", e.ID).Scan(&e.NotGoingCount)
		db.QueryRow("SELECT response FROM event_responses WHERE event_id = ? AND user_id = ?", e.ID, userID).Scan(&e.UserResponse)

		events = append(events, e)
	}

	return events, nil
}

// CreateGroupEvent inserts a new group event into the database
func CreateGroupEvent(db *sql.DB, groupID int, title, description, eventDate, userResponse string, userID int) (int,error) {
	formattedDate, err := utils.ValidateAndFormatDate(eventDate)
	if err != nil {
		return 0,err
	}

	var eventID int
	err = db.QueryRow("INSERT INTO group_events (group_id, title, description, event_date, created_id) VALUES (?, ?, ?, ?, ?) RETURNING id",groupID, title, description, formattedDate, userID,).Scan(&eventID)
	if err != nil {
		return 0,err
	}
	_, err = db.Exec("INSERT INTO event_responses (event_id, user_id, response) VALUES (?, ?, ?)", eventID, userID, userResponse)
	return eventID,err
}

// UpdateUserEventResponse updates or deletes a user's response to an event
func UpdateUserEventResponse(db *sql.DB, eventID int64, userID int, response string) error {
	var existingResponse string
	db.QueryRow("SELECT response FROM event_responses WHERE event_id = ? AND user_id = ?", eventID, userID).Scan(&existingResponse)

	if existingResponse == response {
		_, err := db.Exec("DELETE FROM event_responses WHERE event_id = ? AND user_id = ?", eventID, userID)
		return err
	} else if existingResponse != "" {
		_, err := db.Exec("UPDATE event_responses SET response = ? WHERE event_id = ? AND user_id = ?", response, eventID, userID)
		return err
	} else {
		_, err := db.Exec("INSERT INTO event_responses (event_id, user_id, response) VALUES (?, ?, ?)", eventID, userID, response)
		return err
	}
}
