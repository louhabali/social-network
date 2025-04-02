package handlers

import (
	"encoding/json"
	"main/pkg/app/models"
	"net/http"
	"strconv"
)

func (h *Handler) InviteUsersGroup(w http.ResponseWriter, r *http.Request) {

	groupID := r.PathValue("groupId")
	search := r.URL.Query().Get("search")
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		http.Error(w, "Invalid offset", http.StatusBadRequest)
		return
	}
	users, err := models.GetUsersForGroupInvitation(h.db, groupID, search, offset)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
