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

type AppRepository struct {
	db repository.DBExecutor
}

func NewAppRepository(db repository.DBExecutor) *AppRepository {
	return &AppRepository{db: db}
}

func (r *AppRepository) Create(ctx context.Context, app *domain.App) error {
	query := `
		INSERT INTO apps (id, tenant_id, name, app_code, vapid_public_key, vapid_private_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query, app.ID, app.TenantID, app.Name, app.AppCode, app.VAPIDPublicKey, app.VAPIDPrivateKey, app.CreatedAt, app.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create app: %w", err)
	}
	return nil
}

func (r *AppRepository) GetByCode(ctx context.Context, appCode string) (*domain.App, error) {
	query := `
		SELECT id, tenant_id, name, app_code, vapid_public_key, vapid_private_key, created_at, updated_at
		FROM apps WHERE app_code = $1
	`
	var app domain.App
	err := r.db.QueryRowContext(ctx, query, appCode).Scan(
		&app.ID, &app.TenantID, &app.Name, &app.AppCode, &app.VAPIDPublicKey, &app.VAPIDPrivateKey, &app.CreatedAt, &app.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get app by code: %w", err)
	}
	return &app, nil
}

func (r *AppRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.App, error) {
	query := `
		SELECT id, tenant_id, name, app_code, vapid_public_key, vapid_private_key, created_at, updated_at
		FROM apps WHERE id = $1
	`
	var app domain.App
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&app.ID, &app.TenantID, &app.Name, &app.AppCode, &app.VAPIDPublicKey, &app.VAPIDPrivateKey, &app.CreatedAt, &app.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get app by id: %w", err)
	}
	return &app, nil
}

func (r *AppRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*domain.App, error) {
	query := `
		SELECT id, tenant_id, name, app_code, vapid_public_key, vapid_private_key, created_at, updated_at
		FROM apps WHERE tenant_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list apps: %w", err)
	}
	defer rows.Close()

	var apps []*domain.App
	for rows.Next() {
		var app domain.App
		if err := rows.Scan(&app.ID, &app.TenantID, &app.Name, &app.AppCode, &app.VAPIDPublicKey, &app.VAPIDPrivateKey, &app.CreatedAt, &app.UpdatedAt); err != nil {
			return nil, err
		}
		apps = append(apps, &app)
	}
	return apps, nil
}
