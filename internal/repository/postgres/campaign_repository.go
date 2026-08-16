package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/andrey57x/pwa-webpush-saas/internal/repository"
	"github.com/google/uuid"
)

type CampaignRepository struct {
	db repository.DBExecutor
}

func NewCampaignRepository(db repository.DBExecutor) *CampaignRepository {
	return &CampaignRepository{db: db}
}

func (r *CampaignRepository) Create(ctx context.Context, campaign *domain.Campaign) error {
	query := `
		INSERT INTO campaigns (id, app_id, created_by_user_id, title, body, icon_url, target_url, status, total_targeted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		campaign.ID, campaign.AppID, campaign.CreatedByUserID, campaign.Title, campaign.Body,
		campaign.IconURL, campaign.TargetURL, campaign.Status, campaign.TotalTargeted, campaign.CreatedAt, campaign.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create campaign: %w", err)
	}
	return nil
}

func (r *CampaignRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Campaign, error) {
	query := `
		SELECT id, app_id, created_by_user_id, title, body, icon_url, target_url, status, total_targeted, created_at, updated_at
		FROM campaigns WHERE id = $1
	`
	var c domain.Campaign
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.AppID, &c.CreatedByUserID, &c.Title, &c.Body, &c.IconURL, &c.TargetURL,
		&c.Status, &c.TotalTargeted, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get campaign by id: %w", err)
	}
	return &c, nil
}

func (r *CampaignRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.CampaignStatus, totalTargeted int) error {
	query := `
		UPDATE campaigns 
		SET status = $1, total_targeted = $2, updated_at = NOW() 
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, status, totalTargeted, id)
	if err != nil {
		return fmt.Errorf("failed to update campaign status: %w", err)
	}
	return nil
}

func (r *CampaignRepository) ListByAppID(ctx context.Context, appID uuid.UUID) ([]*domain.Campaign, error) {
	query := `
		SELECT id, app_id, created_by_user_id, title, body, icon_url, target_url, status, total_targeted, created_at, updated_at
		FROM campaigns WHERE app_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to list campaigns: %w", err)
	}
	defer rows.Close()

	var campaigns []*domain.Campaign
	for rows.Next() {
		var c domain.Campaign
		if err := rows.Scan(&c.ID, &c.AppID, &c.CreatedByUserID, &c.Title, &c.Body, &c.IconURL, &c.TargetURL, &c.Status, &c.TotalTargeted, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, &c)
	}
	return campaigns, nil
}
