package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/DKeshavarz/Ar-messenger/internal/models"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

var counter = 0

type RedpandaMessageRepository struct {
	client      *kgo.Client
	adminClient *kadm.Client
	brokers     []string
}

func NewRedpandaMessageRepository(brokers []string) (*RedpandaMessageRepository, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.AllowAutoTopicCreation(),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redpanda client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to Redpanda: %w", err)
	}

	return &RedpandaMessageRepository{
		client:      client,
		adminClient: kadm.NewClient(client),
		brokers:     brokers,
	}, nil
}

func (r *RedpandaMessageRepository) PublishMessage(ctx context.Context, chatID string, msg models.Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	topic := "room-" + chatID
	if err := r.ensureTopic(ctx, topic); err != nil {
		return fmt.Errorf("failed to ensure topic exists: %w", err)
	}

	errCh := make(chan error, 1)
	r.client.Produce(ctx, &kgo.Record{
		Topic: topic,
		Value: b,
		Key:   []byte(msg.RoomName), // For consistent partitioning
	}, func(_ *kgo.Record, err error) {
		if err != nil {
			errCh <- fmt.Errorf("failed to produce message: %w", err)
			return
		}
		errCh <- nil
	})

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *RedpandaMessageRepository) ConsumeMessages(ctx context.Context, chatID string, broadcast chan<- models.Message) error {
	topic := "room-" + chatID

	consumerGroupID := fmt.Sprintf("chatapp-group-%s", chatID)
	client, err := kgo.NewClient(
		kgo.SeedBrokers(r.brokers...),
		kgo.ConsumerGroup(consumerGroupID),
		kgo.ConsumeTopics(topic),
		kgo.BlockRebalanceOnPoll(),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		return fmt.Errorf("failed to create consumer client: %w", err)
	}
	defer client.Close()

	fmt.Printf("Starting consumer for room %s with group %s\n", chatID, consumerGroupID)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			fetches := client.PollFetches(ctx)
			if fetches.IsClientClosed() {
				return fmt.Errorf("client closed")
			}
			if errs := fetches.Errors(); len(errs) > 0 {
				log.Printf("consumer errors: %v", errs)
				continue
			}
			iter := fetches.RecordIter()
			for !iter.Done() {
				record := iter.Next()
				var msg models.Message
				if err := json.Unmarshal(record.Value, &msg); err != nil {
					fmt.Printf("Error decoding message: %v\n", err)
					continue
				}
				log.Printf("Consumed message: content=%s, room=%s, user=%s, partition=%d, offset=%d\n",
					msg.Content, msg.RoomName, msg.Username, record.Partition, record.Offset)

				broadcast <- msg
			}
		}
	}
}

func (r *RedpandaMessageRepository) GetMessageHistory(ctx context.Context, chatID string, msgs chan<- models.Message) error {
	topic := "room-" + chatID

	client, err := kgo.NewClient(
		kgo.SeedBrokers(r.brokers...),
		kgo.ConsumeTopics(topic),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		return fmt.Errorf("failed to create history client: %w", err)
	}
	defer client.Close()

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	for {
		select {
		case <-timeoutCtx.Done():
			return nil
		default:
			fetches := client.PollFetches(timeoutCtx)
			if fetches.IsClientClosed() {
				return nil
			}
			if errs := fetches.Errors(); len(errs) > 0 {
				log.Printf("History fetch errors: %v", errs)
				continue
			}

			iter := fetches.RecordIter()
			for !iter.Done() {
				record := iter.Next()

				var msg models.Message
				if err := json.Unmarshal(record.Value, &msg); err != nil {
					log.Printf("Error decoding history message: %v", err)
					continue
				}

				msgs <- msg
			}
		}
	}
}

func (r *RedpandaMessageRepository) ensureTopic(ctx context.Context, topic string) error {
	// Check if topic exists
	resp, err := r.adminClient.ListTopics(ctx, topic)
	if err != nil {
		return fmt.Errorf("failed to list topics: %w", err)
	}

	if _, exists := resp[topic]; exists {
		return nil
	}

	// Create topic if it doesn't exist
	_, err = r.adminClient.CreateTopics(ctx, 3, 1, nil, topic)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}

	return nil
}

func (r *RedpandaMessageRepository) Close() {
	r.client.Close()
	r.adminClient.Close()
}
