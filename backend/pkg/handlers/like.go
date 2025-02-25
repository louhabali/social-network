package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) LikeHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Pid      string `json:"pd"`
		Interact int `json:"lk"`
	}
	json.NewDecoder(r.Body).Decode(&data)
	userID, ok := r.Context().Value("userID").(int)
	thepostId,_ := strconv.Atoi(data.Pid)
	fmt.Println(thepostId,userID,data.Interact)
	if !ok {
		fmt.Println("invalid userId")
		return
	}

	query := `
		INSERT INTO postInteractions (UserID, PostID, Interaction)
		VALUES (?, ?, ?)
		ON CONFLICT(UserID, PostID) DO UPDATE SET Interaction=excluded.Interaction
	`

	_, err := h.db.Exec(query, userID,thepostId, data.Interact)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
