package handlers

import (
	"encoding/json"
	"main/pkg/app/models"
	"net/http"
	"strconv"
)

func (h *Handler) GroupEventsHandler(w http.ResponseWriter, r *http.Request) {
	groupID, _ := strconv.Atoi(r.PathValue("groupId"))
	userID := r.Context().Value("userID").(int)
	if r.Method == http.MethodGet {
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

		events, err := models.GetGroupEvents(h.db, groupID, offset, userID)
		if err != nil {
			http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	} else if r.Method == http.MethodPost {
		var e models.GroupEvent
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		id,err := models.CreateGroupEvent(h.db, groupID, e.Title, e.Description, e.Date, e.UserResponse, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		e.ID = id
		if e.UserResponse == "going" {
			e.GoingCount++
		}else {
			e.NotGoingCount++
		}
		// notification
		members, err := models.GetGroupMembers(h.db, groupID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		for _, id := range members {
			if id != userID {
				err = models.InsertNotification(h.db, &models.Notification{
					SenderId:  userID,
					ReceverId: id,
					RelatedId: groupID,
					Type:      "group event",
				})
				if err != nil {
					http.Error(w, "Database error", http.StatusInternalServerError)
					return
				}
			}
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(e)
	}
}

func (h *Handler) CreateEventResponse(w http.ResponseWriter, r *http.Request) {
	type EventResponse struct {
		EventID  int64  `json:"eventId"`
		Response string `json:"response"`
	}

	userID := r.Context().Value("userID").(int)

	var e EventResponse
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := models.UpdateUserEventResponse(h.db, e.EventID, userID, e.Response)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
	}

}
