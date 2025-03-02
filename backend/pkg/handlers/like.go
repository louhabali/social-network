package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) LikeHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Pid      string `json:"postId"`
		Interact int `json:"likeStatus"`
	}
	json.NewDecoder(r.Body).Decode(&data)
	userID, ok := r.Context().Value("userID").(int)
	thepostId,_ := strconv.Atoi(data.Pid)
	fmt.Println(thepostId,data.Pid,userID,data.Interact)
	if !ok {
		fmt.Println("invalid userId")
		return
	}

	query := `
		INSERT INTO post_interactions (user_id, post_id, interaction)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, post_id) DO UPDATE SET interaction=excluded.interaction
	`

	_, err := h.db.Exec(query, userID,thepostId, data.Interact)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
