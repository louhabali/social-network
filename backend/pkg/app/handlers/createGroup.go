package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Group struct {
	ID          int    `json:"id"`
	CreatorId   int    `json:"creator_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	NbMembers   int    `json:"nb_members"`
	Status      string `json:"status"`
}

func (h *Handler) CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			fmt.Println("invalid userId")
			return
		}
		var g Group
		err := json.NewDecoder(r.Body).Decode(&g)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		groupId := 0
		err = h.db.QueryRow("INSERT INTO groups (title, description, creator_id) VALUES (?, ?, ?) RETURNING id", g.Title, g.Description, userID).Scan(&groupId)
		if (err != nil) {
			http.Error(w, "Error creating group", http.StatusInternalServerError)
			return
		}
		_, err = h.db.Exec("INSERT INTO group_members (group_id, user_id, status) VALUES (?, ?, ?)", groupId, userID, "creator")
		if (err != nil) {
			http.Error(w, "Error creating group", http.StatusInternalServerError)
			return
		}

		g.CreatorId = userID
		g.ID = groupId
		g.CreatedAt = time.Now().Format("Jan 2, 2006 at 3:04")
		g.Status = "creator"

		// create related chate
		_, err = CreatGroupChat(h.db, userID, g.Title) 
		if (err != nil) {
			http.Error(w, "Error creating group chat", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(g)

	} else if r.Method == http.MethodGet {
		offset := r.URL.Query().Get("offset")
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			fmt.Println("invalid userId")
			return
		}

		rows, _ := h.db.Query(`
		SELECT g.id, g.title, g.description, g.created_at,
       (SELECT COUNT(*) FROM group_members WHERE group_id = g.id AND (status = 'accepted' OR status = 'creator')) AS nb_members,
	   	COALESCE((SELECT status FROM group_members WHERE group_id = g.id AND user_id = ?), '') AS status FROM groups g ORDER BY g.id DESC LIMIT 4 OFFSET ?`, userID, offset)
		defer rows.Close()

		var groups []Group
		for rows.Next() {
			var g Group
			var timeCreated time.Time
			rows.Scan(&g.ID, &g.Title, &g.Description, &timeCreated, &g.NbMembers, &g.Status)
			g.CreatedAt = timeCreated.Format("Jan 2, 2006 at 3:04")
			groups = append(groups, g)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)
	}
}

func (h *Handler) JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	joinSender := r.Context().Value("userID").(int)
	groupID, _ := strconv.Atoi(r.PathValue("groupId"))
	var userID int
	var requestBody struct {
		UserID int `json:"userId"`
	}

	memberInvitation := true
	json.NewDecoder(r.Body).Decode(&requestBody)
	
	if requestBody.UserID == 0 {
		memberInvitation = false
		userID = joinSender
	}else {
		userID = requestBody.UserID
	}

	if r.Method == "POST" {
		res, err := h.db.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", groupID, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		joinReqId, _ := res.LastInsertId()
		if (memberInvitation) {
			insertNotification(h.db, &notificationToInsert{
				RelatedId: int(joinReqId),
				SenderId: joinSender,
				ReciverId: userID,
				Type: "group join invatation",
			})
		} else {
			query := `SELECT creator_id FROM groups  WHERE id = ?`
			var adminId int
			h.db.QueryRow(query,groupID).Scan(&adminId)
			insertNotification(h.db, &notificationToInsert{
				RelatedId: int(joinReqId),
				SenderId: joinSender,
				ReciverId: adminId,
				Type: "group join request",
			})
		}

		fmt.Println("add user")
	} else if r.Method == "DELETE" {
		_, err := h.db.Exec("DELETE FROM group_members WHERE group_id = ? AND user_id = ?", groupID, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}


func getGroupMembers(db *sql.DB, groupId int) ([]int, error){
	var members []int

    query := `
        SELECT user_id 
        FROM group_members 
        WHERE group_id = ? AND (status = 'approved' OR status = 'creator')
    `

    rows, err := db.Query(query, groupId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var userId int
        if err := rows.Scan(&userId); err != nil {
            return nil, err
        }
        members = append(members, userId)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return members, nil
}

func (h *Handler) JoinAccept(w http.ResponseWriter, r *http.Request) {
	fmt.Println("accept group join")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var notif NotifRes
	err := json.NewDecoder(r.Body).Decode(&notif)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	_, err = h.db.Exec("UPDATE group_members SET status = 'approved' WHERE id = ?", notif.RelatedId)
	if err != nil {
		fmt.Println("update join status",err)
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}

	_, err = h.db.Exec("DELETE FROM notifications WHERE id = ?", notif.Id)
	if err != nil {
		fmt.Println("delet notif",err)
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}

	var chatId int
	joinerId := notif.ReceverId
	
	if (notif.Type == "group join request") {
		joinerId = notif.SenderId
	}
	// fmt.Println(notif.RelatedName)
	h.db.QueryRow("SELECT id FROM chats WHERE name = ?", notif.GroupName).Scan(&chatId)
	fmt.Println(chatId)

	addChatMember(h.db, chatId, joinerId)

	// fmt.Println(err)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) JoinReject(w http.ResponseWriter, r *http.Request) {
	fmt.Println("reject group join")
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var notif NotifRes
	err := json.NewDecoder(r.Body).Decode(&notif)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	_, err = h.db.Exec("DELETE FROM group_members WHERE id = ?", notif.RelatedId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}
	_, err = h.db.Exec("DELETE FROM notifications WHERE id = ?", notif.Id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}

	fmt.Println("reject group",notif.Id, notif.RelatedId)
	w.WriteHeader(http.StatusOK)
}