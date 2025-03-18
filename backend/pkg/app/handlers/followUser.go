package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.PathValue("userId"))
	followerID, ok := r.Context().Value("userID").(int)
	if !ok {
		fmt.Println("invalid userId")
		return
	}

	if r.Method == "POST" {
		res, err := h.db.Exec("INSERT INTO follow_requests (follower_id, following_id) VALUES (?, ?)", followerID, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		reqId, err := res.LastInsertId()
		
		err = insertNotification(h.db, &notificationToInsert{
			SenderId: followerID,
			ReciverId: userID,
			RelatedId: int(reqId),
			Type: "follow reqeust",
		})

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

	} else if r.Method == "DELETE" {
		_, err := h.db.Exec("DELETE FROM follow_requests WHERE follower_id = ? AND following_id = ?", followerID, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		_, err = h.db.Exec("DELETE FROM notifications WHERE sender_id = ? AND recever_id = ?", followerID, userID) 
		if err != nil {
			fmt.Println(err)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AcceptFollowReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var notif NotifRes
	err := json.NewDecoder(r.Body).Decode(&notif)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	// fmt.Println(data)
	_, err = h.db.Exec("UPDATE follow_requests SET status = 'approved' WHERE id = ?", notif.RelatedId)
	if err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	_, err = h.db.Exec("DELETE FROM notifications WHERE id = ?", notif.Id)
	if err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	err = createChat(h.db, notif.RelatedId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error creatint chat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) RejectFollowReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var notif NotifRes
	err := json.NewDecoder(r.Body).Decode(&notif)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	// fmt.Println(notif)
	_, err = h.db.Exec("DELETE FROM follow_requests  WHERE id = ?", notif.RelatedId)
	if err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	_, err = h.db.Exec("DELETE FROM notifications WHERE id = ?", notif.Id)
	if err != nil {
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}
	fmt.Println("ooooooooooooooo")
	w.WriteHeader(http.StatusOK)
}

func createChat(db *sql.DB, followReqId int) error {
	var user1, user2 int
	query := `
		SELECT	fr.follower_id,
				fr.following_id
			FROM follow_requests as fr 
			WHERE fr.id = ?
	`
	row := db.QueryRow(query, followReqId)
	err := row.Scan(&user1, &user2)
	if err != nil {
		return err
	}

	fmt.Println(user1, user2)
	_, err = createPrivateChat(db, user1, user2)

	return err
}
