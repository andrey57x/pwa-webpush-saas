package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/google/uuid"
)

type FeedbackUsecase struct {
	logRepo DeliveryLogRepository
}

func NewFeedbackUsecase(logRepo DeliveryLogRepository) *FeedbackUsecase {
	return &FeedbackUsecase{logRepo: logRepo}
}

type FeedbackPingRequest struct {
	CampaignID     uuid.UUID             `json:"campaign_id"`
	SubscriptionID uuid.UUID             `json:"subscription_id"`
	Status         domain.DeliveryStatus `json:"status"` // DELIVERED или CLICKED
}

func (u *FeedbackUsecase) ProcessFeedback(ctx context.Context, req *FeedbackPingRequest) error {
	now := time.Now()
	log := &domain.DeliveryLog{
		ID:             uuid.New(),
		CampaignID:     req.CampaignID,
		SubscriptionID: req.SubscriptionID,
		Status:         req.Status,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := u.logRepo.BatchInsert(ctx, []*domain.DeliveryLog{log}); err != nil {
		return fmt.Errorf("failed to record feedback log: %w", err)
	}

	return nil
}
