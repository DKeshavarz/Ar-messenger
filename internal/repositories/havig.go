package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/DKeshavarz/Ar-messenger/internal/models"
)

type HavigRepo struct{

}

func NewHavigRepo() *HavigRepo{
	return &HavigRepo{}
}
func (h *HavigRepo)PublishMessage(ctx context.Context, chatID string, msg models.Message) error{
	fmt.Println("push in red pand", msg)
	return  nil
}
func (h *HavigRepo)ConsumeMessages(ctx context.Context, chatID string, broadcast chan<- models.Message) error {
	for {
		time.Sleep(1 * time.Second)
		broadcast <- models.Message{
			RoomName: chatID,
			Content: "test",
			Username: "testuser",
		}
	}
	return  nil
}
