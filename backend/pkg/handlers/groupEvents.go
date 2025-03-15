package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type GroupEvent struct {
	ID          int    `json:"id"`
	GroupID     int    `json:"group_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	UserID      int    `json:"created_id"`
	CreatedAt   string `json:"created_at"`
}

func (h *Handler) GroupEventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		groupID, _ := strconv.Atoi(r.PathValue("groupId"))

		rows, err := h.db.Query("SELECT id, group_id, title, description, event_date, created_id, created_at FROM group_events WHERE group_id = ?", groupID)
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

		_, err := h.db.Exec("INSERT INTO group_events (group_id, title, description, event_date, created_id) VALUES (?, ?, ?, ?, ?)",
			groupID, e.Title, e.Description, e.Date, userID)

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

func (h *Handler) CreateEventResponse(w http.ResponseWriter, r *http.Request) {
	type EventResponse struct {
		EventID  int64  `json:"eventId"`
		Response string `json:"status"`
	}
	
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		fmt.Println("invalid userId")
		return
	}

	var e EventResponse
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var exists bool
	h.db.QueryRow("SELECT 1 FROM event_responses WHERE event_id = ? AND user_id = ?",
		e.EventID, userID).Scan(&exists)

	if exists {
		query := `
			UPDATE event_responses 
			SET response = ?
			WHERE event_id = ? AND user_id = ?
		`
		_, err := h.db.Exec(query, e.Response, e.EventID, userID)
		if err != nil {
			fmt.Println("error updating event response: %w", err)
			return
		}
		return
	}

	query := `
		INSERT INTO event_responses (event_id, user_id, response)
		VALUES (?, ?, ?)
	`
	_, err := h.db.Exec(query, e.EventID, userID, e.Response)
	if err != nil {
		fmt.Println("error creating event response: %w", err)
		return
	}
}
