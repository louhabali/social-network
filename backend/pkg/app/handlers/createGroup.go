package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"main/pkg/app/models"
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
		if g.Title == "" || g.Description == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}
		if len(g.Title) > 50 || len(g.Description) > 200 {
			http.Error(w, "Title or description too long", http.StatusBadRequest)
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
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
		userID := r.Context().Value("userID").(int)

		groups, err := models.GetGroups(h.db, userID, offset)
		if err != nil {
			http.Error(w, "Error retrieving groups", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	joinSender := r.Context().Value("userID").(int)
	groupID, err := strconv.Atoi(r.PathValue("groupId"))
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		UserID int `json:"userId"`
	}

	json.NewDecoder(r.Body).Decode(&requestBody)

	memberInvitation := requestBody.UserID != 0
	userID := requestBody.UserID
	if !memberInvitation {
		userID = joinSender
	}

	if r.Method == "POST" {
		joinReqId, err := models.JoinGroup(h.db, groupID, userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if memberInvitation {
			notif, err := models.InsertNotification(h.db, &models.Notification{
				RelatedId: int(joinReqId),
				SenderId:  joinSender,
				ReceverId: userID,
				Type:      "group join invitation",
			})
			if err != nil {
				http.Error(w, "Error inserting notification", http.StatusInternalServerError)
				return
			}
			wsManager.SendToUsers([]int{notif.ReceverId}, WebSocketMessage{
				Type:    "notification",
				Payload: notif,
			})
		} else {
			adminId, err := models.GetGroupCreator(h.db, groupID)
			if err != nil {
				http.Error(w, "Error retrieving group creator", http.StatusInternalServerError)
				return
			}
			notif, err := models.InsertNotification(h.db, &models.Notification{
				RelatedId: int(joinReqId),
				SenderId:  joinSender,
				ReceverId: adminId,
				Type:      "group join request",
			})
			if err != nil {
				http.Error(w, "Error inserting notification", http.StatusInternalServerError)
				return
			}
			wsManager.SendToUsers([]int{notif.ReceverId}, WebSocketMessage{
				Type:    "notification",
				Payload: notif,
			})
		}

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
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) JoinAccept(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var notif models.Notification
	if err := json.NewDecoder(r.Body).Decode(&notif); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := models.ApproveJoinRequest(h.db, notif.RelatedId); err != nil {
		http.Error(w, "Error approving join request", http.StatusInternalServerError)
		return
	}

	if err := models.DeleteNotification(h.db, notif.Id, 0, 0); err != nil {
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}

	joinerId := notif.ReceverId
	if notif.Type == "group join request" {
		joinerId = notif.SenderId
	}
	fmt.Println(notif)
	chatId, err := models.GetChatIDByGroupName(h.db, *notif.GroupName)
	if err != nil {
		fmt.Println(chatId)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := models.AddChatMember(h.db, chatId, joinerId); err != nil {
		http.Error(w, "Failed to add user to chat", http.StatusInternalServerError)
		return
	}

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
	userToRemove := notif.ReceverId
	if notif.Type == "group join request" {
		userToRemove = notif.SenderId
	}

	err = models.RemoveGroupMember(h.db, *notif.GroupId, userToRemove)
	if err != nil {
		http.Error(w, "Error rejecting request", http.StatusInternalServerError)
		return
	}

	err = models.DeleteNotification(h.db, notif.Id, 0, 0)
	if err != nil {
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
