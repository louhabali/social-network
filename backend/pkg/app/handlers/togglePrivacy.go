package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"main/pkg/app/models"
)

func (h *Handler) TogglePrivacy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	var req struct {
		IsPublic bool `json:"isPublic"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req);err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = models.ToggleUserPrivacy(h.db, userID, req.IsPublic)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
