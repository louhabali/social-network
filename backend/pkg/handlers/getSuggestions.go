package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h Handler) Suggestions(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `
	SELECT users.id, users.first_name, users.last_name, users.avatar
	FROM users
	WHERE users.id != ? AND users.id NOT IN (
		SELECT following_id FROM follow_requests WHERE follower_id = ?
	)
	`
	rows, err := h.db.Query(query, userID, userID)
	if err != nil {
		http.Error(w, "Failed to fetch suggestions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notFollowing []Follower
	for rows.Next() {
		var nf Follower
		rows.Scan(&nf.ID, &nf.Firstname, &nf.Lastname, &nf.Avatar)
		if nf.Avatar != nil {
			*nf.Avatar, _ = EncodeImg(*nf.Avatar)
		}
		notFollowing = append(notFollowing, nf)
	}
	fmt.Println("nf: ", notFollowing)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notFollowing)
}
