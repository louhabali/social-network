package handlers

import (
	"encoding/json"
	"fmt"
	"main/pkg/app/models"
	"net/http"
)

func (h *Handler) GetConversations(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(int)

	conversations, err := models.GetConversations(h.db, userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)

}
