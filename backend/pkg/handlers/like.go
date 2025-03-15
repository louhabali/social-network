package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
<<<<<<< HEAD
	"strconv"
=======
>>>>>>> alouhab
)

func (h *Handler) LikeHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
<<<<<<< HEAD
		Pid      string `json:"postId"`
=======
		Pid      int `json:"postId"`
>>>>>>> alouhab
		Interact int `json:"likeStatus"`
	}
	json.NewDecoder(r.Body).Decode(&data)
	userID, ok := r.Context().Value("userID").(int)
<<<<<<< HEAD
	thepostId,_ := strconv.Atoi(data.Pid)
	fmt.Println(thepostId,data.Pid,userID,data.Interact)
=======

	fmt.Println(data.Pid, userID, data.Interact)
>>>>>>> alouhab
	if !ok {
		fmt.Println("invalid userId")
		return
	}

	query := `
		INSERT INTO post_interactions (user_id, post_id, interaction)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, post_id) DO UPDATE SET interaction=excluded.interaction
	`

<<<<<<< HEAD
	_, err := h.db.Exec(query, userID,thepostId, data.Interact)
=======
	_, err := h.db.Exec(query, userID, data.Pid, data.Interact)
>>>>>>> alouhab
	if err != nil {
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
