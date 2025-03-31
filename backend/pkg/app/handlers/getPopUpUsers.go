package handlers

import (
	"encoding/json"
	"fmt"
	"main/pkg/app/models"
	"net/http"
)

func (h *Handler) GetAllusers(w http.ResponseWriter, r *http.Request) {
	myId, _ := r.Context().Value("userID").(int)
	
	followers ,err := models.GetFollowing(h.db,myId)

	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}
