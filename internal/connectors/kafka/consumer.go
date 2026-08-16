package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

// TaskProcessor — интерфейс обработчика задач (бизнес-логика шифрования и отправки)
type TaskProcessor interface {
	ProcessTask(ctx context.Context, task *PushJobTask) error
}

type Consumer struct {
	reader      *kafka.Reader
	processor   TaskProcessor
	workerCount int
}

func NewConsumer(brokers []string, topic, groupID string, workerCount int, processor TaskProcessor) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10,        // 10B
		MaxBytes:       10e6,       // 10MB
		MaxWait:        500 * time.Millisecond,
		CommitInterval: time.Second, // Авто-коммит со смещением
	})

	log.Printf("Kafka Consumer initialized [Topic: %s, Group: %s, Workers: %d]", topic, groupID, workerCount)
	return &Consumer{
		reader:      reader,
		processor:   processor,
		workerCount: workerCount,
	}
}

// Start запускает цикл чтения и распараллеливание через Worker Pool
func (c *Consumer) Start(ctx context.Context) {
	taskChan := make(chan *PushJobTask, c.workerCount*2)
	var wg sync.WaitGroup

	// 1. Запуск Worker Pool
	for i := 0; i < c.workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-taskChan:
					if !ok {
						return
					}
					if err := c.processor.ProcessTask(ctx, task); err != nil {
						log.Printf("[Worker %d] Error processing push job: %v", workerID, err)
					}
				}
			}
		}(i)
	}

	// 2. Главный цикл чтения сообщений из Kafka
	go func() {
		defer close(taskChan)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := c.reader.ReadMessage(ctx)
				if err != nil {
					if ctx.Err() != nil {
						return
					}
					log.Printf("Error reading message from kafka: %v", err)
					continue
				}

				var task PushJobTask
				if err := json.Unmarshal(msg.Value, &task); err != nil {
					log.Printf("Failed to unmarshal kafka message: %v", err)
					continue
				}

				taskChan <- &task
			}
		}
	}()

	wg.Wait()
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}