package postgres

import (
	"context"
	"fmt"

	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/andrey57x/pwa-webpush-saas/internal/repository"
)

type DeliveryLogRepository struct {
	db repository.DBExecutor
}

func NewDeliveryLogRepository(db repository.DBExecutor) *DeliveryLogRepository {
	return &DeliveryLogRepository{db: db}
}

func (r *DeliveryLogRepository) BatchInsert(ctx context.Context, logs []*domain.DeliveryLog) error {
	if len(logs) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for batch insert: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO delivery_logs (id, campaign_id, subscription_id, status, error_code, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement for batch insert: %w", err)
	}
	defer stmt.Close()

	for _, log := range logs {
		_, err := stmt.ExecContext(ctx, log.ID, log.CampaignID, log.SubscriptionID, log.Status, log.ErrorCode, log.CreatedAt, log.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to exec batch insert log %s: %w", log.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit batch insert transaction: %w", err)
	}

	return nil
}
