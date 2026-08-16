package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/andrey57x/pwa-webpush-saas/internal/pkg/hash"
	"github.com/google/uuid"
)

type APIKeyUsecase struct {
	apiKeyRepo APIKeyRepository
}

func NewAPIKeyUsecase(apiKeyRepo APIKeyRepository) *APIKeyUsecase {
	return &APIKeyUsecase{apiKeyRepo: apiKeyRepo}
}

type CreateAPIKeyResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	RawKey    string    `json:"raw_key"` // Показывается клиенту ОДИН раз при создании!
	CreatedAt time.Time `json:"created_at"`
}

func (u *APIKeyUsecase) CreateKey(ctx context.Context, tenantID uuid.UUID, name string) (*CreateAPIKeyResponse, error) {
	rawKey, keyHash, err := hash.GenerateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	now := time.Now()
	apiKey := &domain.APIKey{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Name:      name,
		KeyHash:   keyHash,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.apiKeyRepo.Create(ctx, apiKey); err != nil {
		return nil, fmt.Errorf("failed to save api key: %w", err)
	}

	return &CreateAPIKeyResponse{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		RawKey:    rawKey,
		CreatedAt: apiKey.CreatedAt,
	}, nil
}

func (u *APIKeyUsecase) ListKeys(ctx context.Context, tenantID uuid.UUID) ([]*domain.APIKey, error) {
	return u.apiKeyRepo.ListByTenant(ctx, tenantID)
}
