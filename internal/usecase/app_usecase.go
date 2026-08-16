package usecase

import (
	"context"
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

// CreateApp автоматически генерирует VAPID-ключи на Go-бэкенде при создании приложения PWA
func (u *AppUsecase) CreateApp(ctx context.Context, req *CreateAppRequest) (*domain.App, error) {
	// 1. Автоматическая генерация пары VAPID-ключей (ECDSA P-256)
	vapidKeys, err := vapid.GenerateKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to generate vapid keys: %w", err)
	}

	now := time.Now()
	app := &domain.App{
		ID:              uuid.New(), // UUIDv4
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
	return u.appRepo.GetByCode(ctx, appCode)
}
