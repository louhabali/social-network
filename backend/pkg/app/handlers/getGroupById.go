package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"main/pkg/app/models"
)

func (h *Handler) GetGroupByIdHandler(w http.ResponseWriter, r *http.Request) {
	groupID, _ := strconv.Atoi(r.PathValue("groupId"))
	userID := r.Context().Value("userID").(int)

	group, err := models.GetGroup(h.db, groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get the user's membership status in the group
	group.Status, err = models.GetUserGroupStatus(h.db, groupID, userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(group)
}
