package handlers

import (
	"encoding/json"
	"main/pkg/app/models"
	"net/http"
)

func (h *Handler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Token string `json:"token"`
	}
	var res Response
	if err := json.NewDecoder(r.Body).Decode(&res);err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	_, err := models.GetUserIDFromToken(res.Token, h.db)

	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Token is valid"})
}
