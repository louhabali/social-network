package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"main/pkg/app/models"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketMessage struct {
	Type    string         `json:"type"`
	Payload any `json:"payload"`
}

type WebSocketManager struct {
	connections map[int][]*websocket.Conn
	tokens      map[int]string
	mu          sync.RWMutex
}

var (
	wsManager = &WebSocketManager{
		connections: make(map[int][]*websocket.Conn),
		tokens:      make(map[int]string),
	}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// basicly this implementation
		// it just a reminder to properly check it in production
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins in development
		},
	}
)

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(int)
	cookie, _ := r.Cookie("Token")

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}

	// Register the new connection
	wsManager.registerConnection(userID, conn, cookie.Value)

	userInfo := WebSocketMessage {
		Type: "userInfo",
		Payload: struct{UserId int `json:"user_id"`}{UserId: userID},
	}
	conn.WriteJSON(userInfo)
	go wsManager.handleMessages(userID, conn, h.db)
}

func (wm *WebSocketManager) registerConnection(userID int, conn *websocket.Conn, token string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// Initialize slice if it doesn't exist
	if _, exists := wm.connections[userID]; !exists {
		wm.connections[userID] = make([]*websocket.Conn, 0)
	}

	wm.connections[userID] = append(wm.connections[userID], conn)
	wm.tokens[userID] = token
}

func (wm *WebSocketManager) handleMessages(userID int, conn *websocket.Conn, db *sql.DB) {
	defer func() {
		wm.removeConnection(userID, conn)
	}()

	for {
		// Read Message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway, // User closed tab/browser
				// Connection lost unexpectedly
				websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			return
		}

		if messageType == websocket.TextMessage {
			var msg WebSocketMessage
			if err := json.Unmarshal(p, &msg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
				// Skip this message, wait for next one. Skips invalid message
				// Goes back to waiting for next message
			}

			// Handle different message types
			if msg.Type == "new_message" {
				wm.handleNewMessage(userID, msg.Payload, db)
			}
		}
	}
}

func (wm *WebSocketManager) removeConnection(userID int, conn *websocket.Conn) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	connections, exists := wm.connections[userID]
	if !exists {
		return
	}

	// Find and remove the specific connection
	for i, c := range connections {
		if c == conn {
			wm.connections[userID] = append(connections[:i], connections[i+1:]...)
			break
		}
	}

	// If no more connections and it's not a reconnection attempt
	if len(wm.connections[userID]) == 0 {
		delete(wm.connections, userID)
		delete(wm.tokens, userID)
	}
}

func (wm *WebSocketManager) handleNewMessage(senderID int, payload interface{}, db *sql.DB) {
	payloadBytes, _ := json.Marshal(payload)
	var messageData struct {
		ChatId  int    `json:"chat_id"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(payloadBytes, &messageData); err != nil {
		log.Printf("Error unmarshaling message payload: %v", err)
		return
	}

	msg, err := models.InsertMessage(db,senderID, messageData.ChatId, messageData.Content)
	if err != nil {
		// ana khasni nmchi daba ,ghi kheliha 7ta n9adoha mn b3d , matnsach tpushi dakchi  
		if err.Error() == "At least one follow the other constraint" {
			notification := WebSocketMessage{
				Type:    "error_message",
				Payload: msg,
			}
			wm.SendToUsers([]int{senderID}, notification)
		}
		return
	}

	// Prepare message notification
	notification := WebSocketMessage{
		Type:    "new_message",
		Payload: msg,
	}

	// fmt.Println("cahtid =>",msg.ChatID, "reciver =>", messageData.ReciverId)

	recivers, _ := models.GetUsersInChat(db,msg.ChatID)
	wm.SendToUsers(recivers, notification)
}

func (wm *WebSocketManager) SendToUsers(recivers []int, msg WebSocketMessage) {
	for _, userID := range recivers {
		wm.mu.RLock()
		connections, exists := wm.connections[userID]
		wm.mu.RUnlock()
		if exists {
			// msg.Payload.IsSender = userID == msg.Payload.SenderID
			for _, conn := range connections {
				if err := conn.WriteJSON(msg); err != nil {
					log.Printf("Error sending message to user %d: %v", userID, err)
					wm.removeConnection(userID, conn)
				}
			}
		}
	}
}