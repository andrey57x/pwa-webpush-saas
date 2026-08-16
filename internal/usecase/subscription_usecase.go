package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/google/uuid"
)

type SubscriptionUsecase struct {
	subRepo SubscriptionRepository
	appRepo AppRepository
}

func NewSubscriptionUsecase(subRepo SubscriptionRepository, appRepo AppRepository) *SubscriptionUsecase {
	return &SubscriptionUsecase{
		subRepo: subRepo,
		appRepo: appRepo,
	}
}

type SubscribeRequest struct {
	AppCode        string  `json:"app_code"`
	UserIdentifier *string `json:"user_identifier,omitempty"`
	Endpoint       string  `json:"endpoint"`
	P256dhKey      string  `json:"p256dh_key"`
	AuthKey        string  `json:"auth_key"`
	Browser        string  `json:"browser"`
	OS             string  `json:"os"`
}

func (u *SubscriptionUsecase) RegisterSubscription(ctx context.Context, req *SubscribeRequest) (*domain.Subscription, error) {
	app, err := u.appRepo.GetByCode(ctx, req.AppCode)
	if err != nil || app == nil {
		return nil, fmt.Errorf("app not found for code [%s]", req.AppCode)
	}

	browser := req.Browser

	osInfo := req.OS

	now := time.Now()
	sub := &domain.Subscription{
		ID:             uuid.New(),
		AppID:          app.ID,
		UserIdentifier: req.UserIdentifier,
		Endpoint:       req.Endpoint,
		P256dhKey:      req.P256dhKey,
		AuthKey:        req.AuthKey,
		Browser:        browser,
		OS:             osInfo,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := u.subRepo.Create(ctx, sub); err != nil {
		return nil, fmt.Errorf("failed to save subscription: %w", err)
	}

	return sub, nil
}
