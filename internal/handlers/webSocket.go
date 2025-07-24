package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/DKeshavarz/Ar-messenger/internal/models"
	"github.com/DKeshavarz/Ar-messenger/internal/transport/ws"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type RoomInterface interface {
	SendMessage(ctx context.Context, chatID string, msg models.Message) error
	GetOrCreateRoom(ctx context.Context, chatID string) *models.Room
}

type WebSocketHandler struct {
	svc RoomInterface
}

func NewWebSocketHandler(svc RoomInterface) *WebSocketHandler {
	return &WebSocketHandler{svc: svc}
}

// HandleWebSocket handles the WebSocket connection for a chat room.
// It upgrades the HTTP connection to a WebSocket connection and manages
// the communication between clients in the specified chat room.
//
// The function expects the following URL parameters:
// - chatName: The name of the chat room (required).
// - username: The username of the client (required).
//
// If either the chatName or username is missing, it responds with a
// 400 Bad Request error.
//
// Upon successful connection, it creates a new client and joins the
// specified chat room. It listens for incoming messages from the client
// and validates the message content. If the message is valid, it sends
// the message to the chat room. Any errors during reading or sending
// messages are logged.
//
// The function also ensures that the client is removed from the chat room
// when the connection is closed.
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["chatName"]
	username := r.URL.Query().Get("username")
	if chatID == "" || username == "" {
		http.Error(w, "chatName and username required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &models.Client{
		Conn:     conn,
		Username: username,
		RoomName: chatID,
	}

	room := h.svc.GetOrCreateRoom(r.Context(), chatID)
	room.Join <- client

	defer func() {
		room.Leave <- client
	}()

	for {
		var msg ws.Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Read error:", err)
			return
		}
		if msg.Username == "" || msg.Text == "" || msg.ChatID != chatID {
			continue
		}
		if err := h.svc.SendMessage(r.Context(), chatID, models.Message{
			Username: msg.Username,
			Content:  msg.Text,
			RoomName: chatID,
		}); err != nil {
			log.Println("Publish error:", err)
		}
	}
}
