package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/andrey57x/pwa-webpush-saas/internal/repository"
	"github.com/google/uuid"
)

type SubscriptionRepository struct {
	db repository.DBExecutor // Чистый Go database/sql интерфейс
}

func NewSubscriptionRepository(db repository.DBExecutor) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *domain.Subscription) error {
	query := `
		INSERT INTO subscriptions (
			id, app_id, user_identifier, endpoint, p256dh_key, auth_key, 
			browser, os, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (endpoint) DO UPDATE SET
			is_active = TRUE,
			deactivated_at = NULL,
			p256dh_key = EXCLUDED.p256dh_key,
			auth_key = EXCLUDED.auth_key,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.ExecContext(ctx, query,
		sub.ID, sub.AppID, sub.UserIdentifier, sub.Endpoint, sub.P256dhKey, sub.AuthKey,
		sub.Browser, sub.OS, sub.IsActive, sub.CreatedAt, sub.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert subscription: %w", err)
	}
	return nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	query := `
		SELECT id, app_id, user_identifier, endpoint, p256dh_key, auth_key, browser, os, is_active, deactivated_at, created_at, updated_at
		FROM subscriptions WHERE id = $1
	`
	var sub domain.Subscription
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sub.ID, &sub.AppID, &sub.UserIdentifier, &sub.Endpoint, &sub.P256dhKey, &sub.AuthKey,
		&sub.Browser, &sub.OS, &sub.IsActive, &sub.DeactivatedAt, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get subscription by id: %w", err)
	}
	return &sub, nil
}

func (r *SubscriptionRepository) GetByEndpoint(ctx context.Context, endpoint string) (*domain.Subscription, error) {
	query := `
		SELECT id, app_id, user_identifier, endpoint, p256dh_key, auth_key, browser, os, is_active, deactivated_at, created_at, updated_at
		FROM subscriptions WHERE endpoint = $1
	`
	var sub domain.Subscription
	err := r.db.QueryRowContext(ctx, query, endpoint).Scan(
		&sub.ID, &sub.AppID, &sub.UserIdentifier, &sub.Endpoint, &sub.P256dhKey, &sub.AuthKey,
		&sub.Browser, &sub.OS, &sub.IsActive, &sub.DeactivatedAt, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get subscription by endpoint: %w", err)
	}
	return &sub, nil
}

func (r *SubscriptionRepository) GetActiveByAppID(ctx context.Context, appID uuid.UUID, limit, offset int) ([]*domain.Subscription, error) {
	query := `
		SELECT id, app_id, user_identifier, endpoint, p256dh_key, auth_key, browser, os, is_active, deactivated_at, created_at, updated_at
		FROM subscriptions
		WHERE app_id = $1 AND is_active = TRUE
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, appID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query active subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []*domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		if err := rows.Scan(
			&sub.ID, &sub.AppID, &sub.UserIdentifier, &sub.Endpoint, &sub.P256dhKey, &sub.AuthKey,
			&sub.Browser, &sub.OS, &sub.IsActive, &sub.DeactivatedAt, &sub.CreatedAt, &sub.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subs = append(subs, &sub)
	}
	return subs, nil
}

func (r *SubscriptionRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	query := `
		UPDATE subscriptions 
		SET is_active = FALSE, deactivated_at = $1, updated_at = $1 
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("failed to deactivate subscription: %w", err)
	}
	return nil
}
