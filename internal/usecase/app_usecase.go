package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/andrey57x/pwa-webpush-saas/internal/pkg/crypto/vapid"
	"github.com/google/uuid"
)

type AppUsecase struct {
	appRepo AppRepository
}

func NewAppUsecase(appRepo AppRepository) *AppUsecase {
	return &AppUsecase{appRepo: appRepo}
}

type CreateAppRequest struct {
	TenantID uuid.UUID `json:"tenant_id"`
	Name     string    `json:"name"`
	AppCode  string    `json:"app_code"`
}

func (u *AppUsecase) CreateApp(ctx context.Context, req *CreateAppRequest) (*domain.App, error) {
	vapidKeys, err := vapid.GenerateKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to generate vapid keys: %w", err)
	}

	now := time.Now()
	app := &domain.App{
		ID:              uuid.New(),
		TenantID:        req.TenantID,
		Name:            req.Name,
		AppCode:         req.AppCode,
		VAPIDPublicKey:  vapidKeys.PublicKey,
		VAPIDPrivateKey: vapidKeys.PrivateKey,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := u.appRepo.Create(ctx, app); err != nil {
		return nil, fmt.Errorf("failed to save app: %w", err)
	}

	return app, nil
}

func (u *AppUsecase) GetByCode(ctx context.Context, appCode string) (*domain.App, error) {
	app, err := u.appRepo.GetByCode(ctx, appCode)
	if err != nil {
		return nil, fmt.Errorf("failed to query app: %w", err)
	}
	if app == nil {
		return nil, errors.New("app not found")
	}
	return app, nil
}

func (u *AppUsecase) ListApps(ctx context.Context, tenantID uuid.UUID) ([]*domain.App, error) {
	return u.appRepo.ListByTenant(ctx, tenantID)
}
