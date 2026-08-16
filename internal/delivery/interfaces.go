package delivery

import (
	"context"

	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/andrey57x/pwa-webpush-saas/internal/usecase"
	"github.com/google/uuid"
)

type SubscriptionUsecase interface {
	RegisterSubscription(ctx context.Context, req *usecase.SubscribeRequest) (*domain.Subscription, error)
}

type AppUsecase interface {
	CreateApp(ctx context.Context, req *usecase.CreateAppRequest) (*domain.App, error)
	GetByCode(ctx context.Context, appCode string) (*domain.App, error)
	ListApps(ctx context.Context, tenantID uuid.UUID) ([]*domain.App, error)
}

type CampaignUsecase interface {
	LaunchCampaign(ctx context.Context, tenantID, campaignID uuid.UUID) error
	ListCampaigns(ctx context.Context, tenantID, appID uuid.UUID) ([]*domain.Campaign, error)
}

type FeedbackUsecase interface {
	ProcessFeedback(ctx context.Context, req *usecase.FeedbackPingRequest) error
}

type AuthUsecase interface {
	Register(ctx context.Context, req *usecase.RegisterRequest) (*domain.User, error)
	Login(ctx context.Context, req *usecase.LoginRequest) (*usecase.AuthResponse, error)
}

type APIKeyUsecase interface {
	CreateKey(ctx context.Context, tenantID uuid.UUID, name string) (*usecase.CreateAPIKeyResponse, error)
	ListKeys(ctx context.Context, tenantID uuid.UUID) ([]*domain.APIKey, error)
}

type AnalyticsUsecase interface {
	GetCampaignStats(ctx context.Context, tenantID, campaignID uuid.UUID) (*domain.CampaignStats, error)
}
