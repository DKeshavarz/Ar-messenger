package services

import (
	"context"
	"errors"
	"sync"

	"github.com/DKeshavarz/Ar-messenger/internal/models"
)

type MessageRepository interface {
    PublishMessage(ctx context.Context, chatID string, msg models.Message) error
    ConsumeMessages(ctx context.Context, chatID string, broadcast chan<- models.Message) error
    GetMessageHistory(ctx context.Context, chatID string, msgs chan<- models.Message) error
}

type RoomService struct {
    repo  MessageRepository
    rooms map[string]*models.Room
    mu    sync.RWMutex
}

func NewRoom(chatID string) *models.Room {
    return &models.Room{
		RoomName:  chatID,
        Clients:   make(map[*models.Client]bool),
        Broadcast: make(chan models.Message, 64),
        Join:      make(chan *models.Client),
        Leave:     make(chan *models.Client),
    }
} 

func NewRoomService(repo MessageRepository) *RoomService {
    return &RoomService{
        repo:  repo,
        rooms: make(map[string]*models.Room),
    }
}

func (s *RoomService) SendMessage(ctx context.Context, chatID string, msg models.Message) error {
    if msg.RoomName != chatID || msg.Username == "" || msg.Content == "" {
        return errors.New("invalid message")
    }
    return s.repo.PublishMessage(ctx, chatID, msg)
}

func (s *RoomService) GetOrCreateRoom(ctx context.Context, chatID string) *models.Room {
    s.mu.Lock()
    defer s.mu.Unlock()
    if room, exists := s.rooms[chatID]; exists {
        return room
    }
    room := NewRoom(chatID)
    s.rooms[chatID] = room
    go room.Run()
    go s.repo.ConsumeMessages(ctx, chatID, room.Broadcast)
    return room
}