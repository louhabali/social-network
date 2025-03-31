package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"main/pkg/app/models"
)

func (h *Handler) UserData(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		fmt.Println("invalid userId")
		return
	}

	userInfo, err := models.GetUserInfo(h.db, userID)
	if err != nil {
		fmt.Println("Error fetching profile privacy:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}
