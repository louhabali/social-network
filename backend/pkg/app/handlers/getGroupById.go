package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) GetGroupByIdHandler(w http.ResponseWriter, r *http.Request) {

	groupID, _ := strconv.Atoi(r.PathValue("groupId"))
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		fmt.Println("invalid userId")
		return
	}
	g := GetGroup(groupID, h.db)
	h.db.QueryRow("SELECT status FROM group_members WHERE group_id = ? AND user_id = ?", groupID, userID).Scan(&g.Status)
	if g.CreatorId == userID {
		g.Status = "creator"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

func GetGroup(groupID int, db *sql.DB) Group {
	var g Group
	var timeCreated time.Time

	db.QueryRow("SELECT id, title, description, creator_id, created_at FROM groups WHERE id = ?", groupID).
		Scan(&g.ID, &g.Title, &g.Description, &g.CreatorId, &timeCreated)
	g.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")

	return g
}
