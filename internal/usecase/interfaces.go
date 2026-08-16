package usecase

import (
	"context"
	"time"

	"github.com/andrey57x/pwa-webpush-saas/internal/connectors/kafka"
	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/andrey57x/pwa-webpush-saas/internal/pkg/pushsender"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	GetByEndpoint(ctx context.Context, endpoint string) (*domain.Subscription, error)
	GetActiveByAppID(ctx context.Context, appID uuid.UUID, limit, offset int) ([]*domain.Subscription, error)
	Deactivate(ctx context.Context, id uuid.UUID) error
}

type CampaignRepository interface {
	Create(ctx context.Context, campaign *domain.Campaign) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Campaign, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.CampaignStatus, totalTargeted int) error
}

type DeliveryLogRepository interface {
	BatchInsert(ctx context.Context, logs []*domain.DeliveryLog) error
}

type AppRepository interface {
	Create(ctx context.Context, app *domain.App) error
	GetByCode(ctx context.Context, appCode string) (*domain.App, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.App, error)
}

type DedupRepository interface {
	FilterDuplicates(ctx context.Context, campaignID uuid.UUID, subIDs []uuid.UUID, ttl time.Duration) ([]uuid.UUID, error)
}

type PushJobProducer interface {
	PublishPushJobs(ctx context.Context, tasks []*kafka.PushJobTask) error
}

// PushSender — приведен к 12 аргументам (включая campaignID и subscriptionID)
type PushSender interface {
	SendPush(
		ctx context.Context,
		endpoint, p256dh, auth, title, body, iconURL, targetURL, vapidPublicKey, vapidPrivateKey, campaignID, subscriptionID string,
	) *pushsender.PushResult
}
