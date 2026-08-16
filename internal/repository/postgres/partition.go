package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/andrey57x/pwa-webpush-saas/internal/repository"
)

// EnsureWeeklyPartitions проверяет и создает недельные партиции в БД
func EnsureWeeklyPartitions(ctx context.Context, db repository.DBExecutor) error {
	now := time.Now()

	// Создаем партиции от прошлой недели (-1) до +4 недель вперед
	for i := -1; i <= 4; i++ {
		targetDate := now.AddDate(0, 0, i*7).Format("2006-01-02")
		query := "SELECT create_delivery_logs_partition_for_date($1::DATE);"

		if _, err := db.ExecContext(ctx, query, targetDate); err != nil {
			return fmt.Errorf("failed to auto-create weekly partition for date [%s]: %w", targetDate, err)
		}
	}

	return nil
}

// StartPartitionScheduler запускает проверку при старте И каждый день в фоне раз в 24 часа
func StartPartitionScheduler(ctx context.Context, db repository.DBExecutor) {
	if err := EnsureWeeklyPartitions(ctx, db); err != nil {
		log.Printf("[PartitionScheduler] Error checking weekly partitions: %v", err)
	} else {
		log.Println("[PartitionScheduler] Weekly database partitions verified and ready!")
	}

	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("[PartitionScheduler] Scheduler stopped.")
				return
			case <-ticker.C:
				if err := EnsureWeeklyPartitions(ctx, db); err != nil {
					log.Printf("[PartitionScheduler] Daily background partition creation error: %v", err)
				} else {
					log.Println("[PartitionScheduler] Daily background partition maintenance completed.")
				}
			}
		}
	}()
}
