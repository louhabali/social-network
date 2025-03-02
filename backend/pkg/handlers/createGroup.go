package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Group struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	NbMembers   int    `json:"nb_members"`
}

func (h *Handler) CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			fmt.Println("invalid userId")
			return
		}
		var g Group
		err := json.NewDecoder(r.Body).Decode(&g)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		groupId := 0
		h.db.QueryRow("INSERT INTO groups (title, description, creator_id) VALUES (?, ?, ?) RETURNING id", g.Title, g.Description, userID).Scan(&groupId)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GetGroup(groupId, h.db))

	} else if r.Method == http.MethodGet {
		offset := r.URL.Query().Get("offset")

		rows, _ := h.db.Query("SELECT id, title, description, created_at FROM groups ORDER BY id DESC LIMIT 4 OFFSET ?", offset)
		defer rows.Close()

		var groups []Group
		for rows.Next() {
			var g Group
			var timeCreated time.Time
			rows.Scan(&g.ID, &g.Title, &g.Description, &timeCreated)
			g.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")
			h.db.QueryRow("SELECT COUNT(*) FROM group_members WHERE group_id = ? AND status = 'accepted' ", g.ID).Scan(&g.NbMembers)
			groups = append(groups, g)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)
	}
}
