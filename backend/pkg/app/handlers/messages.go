package handlers

import (
	"encoding/json"
	"main/pkg/app/models"
	"net/http"
	"strconv"
)

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	userID, _ := r.Context().Value("userID").(int)
	chatID := r.URL.Query().Get("chat_id")

	if chatID == "" {
		http.Error(w, "chat_id is required", http.StatusBadRequest)
		return
	}

	chatId, _ := strconv.Atoi(chatID)
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	messages, err := models.GetMessages(h.db, userID, chatId, offset)
	if err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	// fmt.Println(messages)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *Handler) UpdateLastReadMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	userID, _ := r.Context().Value("userID").(int)

	var ulmr struct {
		ChatID int `json:"chat_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&ulmr); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := models.UpdateLastReadMessage(h.db, userID, ulmr.ChatID)
	if err != nil {
		http.Error(w, "Failed to update last read message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
