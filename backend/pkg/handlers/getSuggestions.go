package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) Suggestions(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	rows, _ := h.db.Query(`
		SELECT users.id, users.first_name, users.last_name FROM users
		JOIN follow_requests ON users.id = follow_requests.following_id
		WHERE follow_requests.follower_id = ? AND follow_requests.status != 'approved';
	`, userID)

	defer rows.Close()

	var notFollowing []Follower
	for rows.Next() {
		var nf Follower
		rows.Scan(&nf.ID, &nf.Firstname, &nf.Lastname)
		notFollowing = append(notFollowing, nf)
	}
	fmt.Println("nf: ", notFollowing)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notFollowing)

}
