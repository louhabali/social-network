package handlers

import (
	"encoding/json"
	"fmt"
	"main/pkg/app/models"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userID := r.Context().Value("userID").(int)

		var g models.Group
		err := json.NewDecoder(r.Body).Decode(&g)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		groupId, err := models.CreateGroup(h.db, g.Title, g.Description, userID)
		if err != nil {
			http.Error(w, "Error creating group", http.StatusInternalServerError)
			return
		}

		err = models.AddGroupMember(h.db, groupId, userID, "creator")
		if err != nil {
			http.Error(w, "Error adding creator to group", http.StatusInternalServerError)
			return
		}

		g.CreatorId = userID
		g.ID = groupId
		g.CreatedAt = time.Now().Format("Jan 2, 2006 at 3:04")
		g.Status = "creator"

		_, err = models.CreatGroupChat(h.db, userID, g.Title)
		if err != nil {
			http.Error(w, "Error creating group chat", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(g)

	} else if r.Method == http.MethodGet {
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		userID := r.Context().Value("userID").(int)

		groups, err := models.GetGroups(h.db, userID, offset)
		if err != nil {
			http.Error(w, "Error retrieving groups", http.StatusInternalServerError)
			return
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
	} else {
		userID = requestBody.UserID
	}

	if r.Method == "POST" {
		joinReqId, err := models.JoinGroup(h.db, groupID, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if memberInvitation {
			models.InsertNotification(h.db, &models.Notification{
				RelatedId: int(joinReqId),
				SenderId:  joinSender,
				ReceverId: userID,
				Type:      "group join invatation",
			})
		} else {
			query := `SELECT creator_id FROM groups  WHERE id = ?`
			var adminId int
			h.db.QueryRow(query, groupID).Scan(&adminId)
			models.InsertNotification(h.db, &models.Notification{
				RelatedId: int(joinReqId),
				SenderId:  joinSender,
				ReceverId: adminId,
				Type:      "group join request",
			})
		}

		fmt.Println("add user")
	} else if r.Method == "DELETE" {
		err := models.DeleteGroupMemberNotifications(h.db, groupID, userID)
		if err != nil {
			http.Error(w, "Failed to delete notifications", http.StatusInternalServerError)
			return
		}

		err = models.RemoveGroupMember(h.db, groupID, userID)
		if err != nil {
			http.Error(w, "Error removing user from group", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) JoinAccept(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var notif models.Notification
	err := json.NewDecoder(r.Body).Decode(&notif)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	err = models.ApproveJoinRequest(h.db, notif.RelatedId)
	if err != nil {
		http.Error(w, "Error approving join request", http.StatusInternalServerError)
		return
	}

	err = models.DeleteNotification(h.db, notif.Id)
	if err != nil {
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}

	var chatId int
	joinerId := notif.ReceverId

	if notif.Type == "group join request" {
		joinerId = notif.SenderId
	}
	h.db.QueryRow("SELECT id FROM chats WHERE name = ?", notif.GroupName).Scan(&chatId)
	fmt.Println(chatId)

	models.AddChatMember(h.db, chatId, joinerId)

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) JoinReject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var notif models.Notification
	err := json.NewDecoder(r.Body).Decode(&notif)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	err = models.RemoveGroupMember(h.db, *notif.GroupId, notif.SenderId)
	if err != nil {
		http.Error(w, "Error rejecting request", http.StatusInternalServerError)
		return
	}

	err = models.DeleteNotification(h.db, notif.Id)
	if err != nil {
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}

	fmt.Println("reject group", notif.Id, notif.RelatedId)
	w.WriteHeader(http.StatusOK)
}
