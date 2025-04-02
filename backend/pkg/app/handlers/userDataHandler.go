package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"main/pkg/app/models"
)

func (h *Handler) UserData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	userID := r.Context().Value("userID").(int)
	
	userInfo, err := models.GetUserInfo(h.db, userID)
	if err != nil {
		fmt.Println("Error fetching profile privacy:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}
