package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type folowNotif struct {
	Type      string `json:"type"`
	RequestID int    `json:"request_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userId, _ := r.Context().Value("userID").(int)
	fmt.Println(userId)
	followReqs, err := getFolowReq(h.db, userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching follow requests: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(followReqs); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func getFolowReq(db *sql.DB, userId int) ([]folowNotif, error) {
	fmt.Println("here")
	query := `
	SELECT fr.id AS request_id,
		   u.first_name,
		   u.last_name
	FROM follow_requests fr
	JOIN users u ON fr.follower_id = u.id
	WHERE fr.following_id = ? AND fr.status = 'pending';
`

	rows, err := db.Query(query, userId)
	if err != nil {
		log.Printf("Error fetching follow requests: %v", err)
		return nil, err
	}
	defer rows.Close()

	var requests []folowNotif
	for rows.Next() {
		var requestId int
		var firstName, lastName string
		// var followerId int

		if err := rows.Scan(&requestId, &firstName, &lastName); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		notification := folowNotif{
			RequestID: requestId,
			FirstName: firstName,
			LastName:  lastName,
			Type:      "follow_request",
		}

		requests = append(requests, notification)
	}

	return requests, nil
}
