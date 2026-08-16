package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// PushJobTask — структура события, публикуемого в топик Kafka
type PushJobTask struct {
	CampaignID     uuid.UUID `json:"campaign_id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	Endpoint       string    `json:"endpoint"`
	P256dhKey      string    `json:"p256dh_key"`
	AuthKey        string    `json:"auth_key"`
	Title          string    `json:"title"`
	Body           string    `json:"body"`
	IconURL        string    `json:"icon_url,omitempty"`
	TargetURL      string    `json:"target_url,omitempty"`
	VAPIDPublicKey string    `json:"vapid_public_key"`
	VAPIDPrivKey   string    `json:"vapid_private_key"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{}, // Шардирование по Hash(Key) -> ключевая стратегия для дипломного исследования!
		MaxAttempts:  5,
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
		Async:        false, // Синхронное подтверждение для надежности
	}

	log.Printf("Kafka Producer initialized for topic [%s]", topic)
	return &Producer{writer: writer}
}

func (p *Producer) PublishPushJobs(ctx context.Context, tasks []*PushJobTask) error {
	if len(tasks) == 0 {
		return nil
	}

	messages := make([]kafka.Message, len(tasks))
	for i, task := range tasks {
		payloadBytes, err := json.Marshal(task)
		if err != nil {
			return fmt.Errorf("failed to marshal push job task: %w", err)
		}

		messages[i] = kafka.Message{
			// Ключ партиционирования = subscription_id (для равномерной и упорядоченной доставки)
			Key:   []byte(task.SubscriptionID.String()),
			Value: payloadBytes,
			Time:  time.Now(),
		}
	}

	if err := p.writer.WriteMessages(ctx, messages...); err != nil {
		return fmt.Errorf("failed to write messages to kafka: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
