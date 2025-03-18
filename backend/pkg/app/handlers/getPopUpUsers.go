package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type privateUsers struct {
	Id    int    `json:"id"`
	Fname string `json:"fname"`
	Lname string `json:"lname"`
}

func (h *Handler) GetAllusers(w http.ResponseWriter, r *http.Request) {
	myId, _ := r.Context().Value("userID").(int)
	
	flwrs ,err := getFollowing(h.db,myId)
	//fmt.Println("following is : ",flwrs)
	if err != nil {
		fmt.Println("11",err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flwrs)
}
