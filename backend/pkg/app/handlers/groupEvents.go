package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type GroupEvent struct {
	ID            int    `json:"id"`
	GroupID       int    `json:"group_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Date          string `json:"date"`
	UserID        int    `json:"created_id"`
	CreatedAt     string `json:"created_at"`
	GoingCount    int    `json:"going_count"`
	NotGoingCount int    `json:"not_going_count"`
	UserResponse  string `json:"user_response"`
}

func (h *Handler) GroupEventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		groupID, _ := strconv.Atoi(r.PathValue("groupId"))
		offset := r.URL.Query().Get("offset")
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			fmt.Println("invalid userId")
			return
		}
		rows, err := h.db.Query("SELECT id, group_id, title, description, event_date, created_id, created_at FROM group_events WHERE group_id = ? ORDER BY id DESC LIMIT 4 OFFSET ?", groupID, offset)
		if err != nil {
			http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var events []GroupEvent
		for rows.Next() {
			var e GroupEvent
			err := rows.Scan(&e.ID, &e.GroupID, &e.Title, &e.Description, &e.Date, &e.UserID, &e.CreatedAt)
			if err != nil {
				http.Error(w, "Error processing events", http.StatusInternalServerError)
				return
			}

			h.db.QueryRow("SELECT COUNT(*) FROM event_responses WHERE event_id = ? AND response = 'going'", e.ID).Scan(&e.GoingCount)
			h.db.QueryRow("SELECT COUNT(*) FROM event_responses WHERE event_id = ? AND response = 'not_going'", e.ID).Scan(&e.NotGoingCount)
			h.db.QueryRow("SELECT response FROM event_responses WHERE event_id = ? AND user_id = ?", e.ID, userID).Scan(&e.UserResponse)
			events = append(events, e)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	} else if r.Method == http.MethodPost {
		groupID, _ := strconv.Atoi(r.PathValue("groupId"))
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			fmt.Println("invalid userId")
			return
		}

		var e GroupEvent
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		formattedDate, err := validateAndFormatDate(e.Date)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = h.db.Exec("INSERT INTO group_events (group_id, title, description, event_date, created_id) VALUES (?, ?, ?, ?, ?)",
			groupID, e.Title, e.Description, formattedDate, userID)

		if err != nil {
			http.Error(w, "Failed to create event", http.StatusInternalServerError)
			return
		}

		// eventId, _ := res.LastInsertId();

		// notification
		members, err := getGroupMembers(h.db, groupID)
		if err != nil {
			log.Println("error geting G members for event notif")
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		
		for _ , id := range members {
			if id != userID {
				err = insertNotification(h.db, &notificationToInsert{
					SenderId: userID,
					ReciverId: id,
					RelatedId: groupID,
					Type: "group event",
				})
				if err != nil {
					log.Println("notif insert error")
					http.Error(w, "Database error", http.StatusInternalServerError)
					return
				}
			}
		}

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(e)
	}
}

func validateAndFormatDate(inputDate string) (string, error) {
	parsedTime, err := time.Parse("2006-01-02T15:04", inputDate)
	if err != nil {
		return "", err
	}

	if parsedTime.Before(time.Now()) {
		return "", errors.New("the event date cannot be in the past")
	}

	formattedDate := parsedTime.Format("2006-01-02 15:04:05")
	return formattedDate, nil
}

func (h *Handler) CreateEventResponse(w http.ResponseWriter, r *http.Request) {
	type EventResponse struct {
		EventID  int64  `json:"eventId"`
		Response string `json:"response"`
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var e EventResponse
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var existingResponse string
	h.db.QueryRow("SELECT response FROM event_responses WHERE event_id = ? AND user_id = ?", e.EventID, userID).Scan(&existingResponse)
	fmt.Println(existingResponse, e.Response)
	if existingResponse == e.Response {
		h.db.Exec("DELETE FROM event_responses WHERE event_id = ? AND user_id = ?", e.EventID, userID)
	} else if existingResponse != "" {
		h.db.Exec("UPDATE event_responses SET response = ? WHERE event_id = ? AND user_id = ?", e.Response, e.EventID, userID)
	} else {
		h.db.Exec("INSERT INTO event_responses (event_id, user_id, response) VALUES (?, ?, ?)", e.EventID, userID, e.Response)
	}

}
