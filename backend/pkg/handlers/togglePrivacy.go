package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *Handler) TogglePrivacy(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.PathValue("userId"))

	var req struct {
		IsPublic bool `json:"isPublic"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	_, err := h.db.Exec("UPDATE users SET is_public = ? WHERE id = ?", req.IsPublic, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
