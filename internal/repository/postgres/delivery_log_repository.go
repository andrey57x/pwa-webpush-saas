package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/andrey57x/pwa-webpush-saas/internal/repository"
)

type DeliveryLogRepository struct {
	db repository.DBExecutor
}

func NewDeliveryLogRepository(db repository.DBExecutor) *DeliveryLogRepository {
	return &DeliveryLogRepository{db: db}
}

// BatchInsert выполняет мгновенный массовый INSERT одном SQL-запросом
func (r *DeliveryLogRepository) BatchInsert(ctx context.Context, logs []*domain.DeliveryLog) error {
	if len(logs) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(logs))
	valueArgs := make([]interface{}, 0, len(logs)*7)

	for i, log := range logs {
		n := i * 7
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", n+1, n+2, n+3, n+4, n+5, n+6, n+7))
		valueArgs = append(valueArgs, log.ID, log.CampaignID, log.SubscriptionID, log.Status, log.ErrorCode, log.CreatedAt, log.UpdatedAt)
	}

	query := fmt.Sprintf(`
		INSERT INTO delivery_logs (id, campaign_id, subscription_id, status, error_code, created_at, updated_at)
		VALUES %s
	`, strings.Join(valueStrings, ","))

	_, err := r.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to execute batch insert: %w", err)
	}

	return nil
}
