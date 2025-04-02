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
	userID := r.Context().Value("userID").(int)
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
