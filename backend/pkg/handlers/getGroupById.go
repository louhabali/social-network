package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) GetGroupByIdHandler(w http.ResponseWriter, r *http.Request) {

	groupID, _ := strconv.Atoi(r.PathValue("groupId"))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GetGroup(groupID,h.db))
}


func GetGroup(groupID int, db *sql.DB) Group {
	var g Group
	var timeCreated time.Time

	db.QueryRow("SELECT id, title, description, creator_id, created_at FROM groups WHERE id = ?", groupID).
		Scan(&g.ID, &g.Title, &g.Description, &timeCreated, &g.CreatedAt)
	g.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")

	return g
}