package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"main/pkg/app/models"
)

func (h *Handler) TogglePrivacy(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.PathValue("userId"))

	var req struct {
		IsPublic bool `json:"isPublic"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	err = models.ToggleUserPrivacy(h.db, userID, req.IsPublic)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
