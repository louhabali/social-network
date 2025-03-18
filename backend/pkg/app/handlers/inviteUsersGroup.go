package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handler) InviteUsersGroup(w http.ResponseWriter, r *http.Request) {

	groupID := r.PathValue("groupId")
	search := r.URL.Query().Get("search")
	offset := r.URL.Query().Get("offset")

	query := `
        SELECT id, first_name,last_name FROM users 
        WHERE id NOT IN (SELECT user_id FROM group_members WHERE group_id = ?)
    `
	params := []interface{}{groupID}
	fmt.Println(search)
	if search != "" {
		query += " and first_name like CONCAT('%', ?, '%') OR last_name like CONCAT('%', ?, '%')  OR nickname like CONCAT('%', ?, '%')"
		params = append(params, search,search,search)
	}

	query += " LIMIT ? OFFSET ?"
	params = append(params, 4, offset)

	rows, _ := h.db.Query(query, params...)

	type user struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	var users []user
	for rows.Next() {
		var u user
		var firstname, lastname string
		rows.Scan(&u.ID, &firstname, lastname)
		u.Name = firstname + lastname
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
