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

type APIKeyRepository struct {
	db repository.DBExecutor
}

func NewAPIKeyRepository(db repository.DBExecutor) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

func (r *APIKeyRepository) Create(ctx context.Context, apiKey *domain.APIKey) error {
	query := `
		INSERT INTO api_keys (id, tenant_id, name, key_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query, apiKey.ID, apiKey.TenantID, apiKey.Name, apiKey.KeyHash, apiKey.CreatedAt, apiKey.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create api key: %w", err)
	}
	return nil
}

func (r *APIKeyRepository) GetByHash(ctx context.Context, keyHash string) (*domain.APIKey, error) {
	query := `
		SELECT id, tenant_id, name, key_hash, last_used_at, created_at, updated_at
		FROM api_keys WHERE key_hash = $1
	`
	var key domain.APIKey
	err := r.db.QueryRowContext(ctx, query, keyHash).Scan(
		&key.ID, &key.TenantID, &key.Name, &key.KeyHash, &key.LastUsedAt, &key.CreatedAt, &key.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query api key by hash: %w", err)
	}
	return &key, nil
}

func (r *APIKeyRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*domain.APIKey, error) {
	query := `
		SELECT id, tenant_id, name, key_hash, last_used_at, created_at, updated_at
		FROM api_keys WHERE tenant_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list api keys: %w", err)
	}
	defer rows.Close()

	var keys []*domain.APIKey
	for rows.Next() {
		var key domain.APIKey
		if err := rows.Scan(&key.ID, &key.TenantID, &key.Name, &key.KeyHash, &key.LastUsedAt, &key.CreatedAt, &key.UpdatedAt); err != nil {
			return nil, err
		}
		keys = append(keys, &key)
	}
	return keys, nil
}
