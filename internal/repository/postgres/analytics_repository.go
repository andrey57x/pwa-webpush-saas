package postgres

import (
	"context"
	"fmt"

	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/andrey57x/pwa-webpush-saas/internal/repository"
	"github.com/google/uuid"
)

type AnalyticsRepository struct {
	db repository.DBExecutor
}

func NewAnalyticsRepository(db repository.DBExecutor) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) GetCampaignStatusCounts(ctx context.Context, campaignID uuid.UUID) (map[domain.DeliveryStatus]int, error) {
	query := `
		SELECT status, COUNT(*)
		FROM delivery_logs
		WHERE campaign_id = $1
		GROUP BY status
	`
	rows, err := r.db.QueryContext(ctx, query, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to query campaign status counts: %w", err)
	}
	defer rows.Close()

	counts := make(map[domain.DeliveryStatus]int)
	for rows.Next() {
		var status domain.DeliveryStatus
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan status count: %w", err)
		}
		counts[status] = count
	}

	return counts, nil
}
