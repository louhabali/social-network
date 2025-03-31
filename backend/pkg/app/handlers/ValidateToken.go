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
	json.NewDecoder(r.Body).Decode(&res)
	
	_, err := models.GetUserIDFromToken(res.Token, h.db)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Token is invalid"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Token is valid"})
}
