package handlers

import (
	"encoding/json"
	"main/pkg/app/models"
	"net/http"
)

func (h Handler) Suggestions(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(int)

	followers := models.NotFollowing(h.db, userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)

}
