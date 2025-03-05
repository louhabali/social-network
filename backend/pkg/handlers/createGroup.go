package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Group struct {
	ID          int    `json:"id"`
	CreatorId   int    `json:"creator_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	NbMembers   int    `json:"nb_members"`
	Status      string `json:"status"`
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
		h.db.Exec("INSERT INTO group_members (group_id, user_id, status) VALUES (?, ?, ?)", groupId, userID, "creator")

		g.CreatorId = userID
		g.ID = groupId
		g.CreatedAt = time.Now().Format("Jan 2, 2006 at 3:04")
		g.Status = "creator"

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(g)

	} else if r.Method == http.MethodGet {
		offset := r.URL.Query().Get("offset")
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			fmt.Println("invalid userId")
			return
		}

		rows, _ := h.db.Query(`
		SELECT g.id, g.title, g.description, g.created_at,
       (SELECT COUNT(*) FROM group_members WHERE group_id = g.id AND (status = 'accepted' OR status = 'creator')) AS nb_members,
	   	COALESCE((SELECT status FROM group_members WHERE group_id = g.id AND user_id = ?), '') AS status FROM groups g ORDER BY g.id DESC LIMIT 4 OFFSET ?`, userID, offset)
		defer rows.Close()

		var groups []Group
		for rows.Next() {
			var g Group
			var timeCreated time.Time
			rows.Scan(&g.ID, &g.Title, &g.Description, &timeCreated, &g.NbMembers, &g.Status)
			g.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")
			groups = append(groups, g)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)
	}
}

func (h *Handler) JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	groupID, _ := strconv.Atoi(r.PathValue("groupId"))
	var userID int
	var requestBody struct {
		UserID int `json:"userId"`
	}
	json.NewDecoder(r.Body).Decode(&requestBody)
	if requestBody.UserID == 0 {
		userID = r.Context().Value("userID").(int)
	}else {
		userID = requestBody.UserID
	}

	if r.Method == "POST" {
		_, err := h.db.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", groupID, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	} else if r.Method == "DELETE" {
		_, err := h.db.Exec("DELETE FROM group_members WHERE group_id = ? AND user_id = ?", groupID, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
