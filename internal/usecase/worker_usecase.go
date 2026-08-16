package usecase

import (
	"context"
	"log"
	"time"

	"github.com/andrey57x/pwa-webpush-saas/internal/connectors/kafka"
	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/google/uuid"
)

type WorkerUsecase struct {
	subRepo SubscriptionRepository
	logRepo DeliveryLogRepository
	sender  PushSender
}

func NewWorkerUsecase(subRepo SubscriptionRepository, logRepo DeliveryLogRepository, sender PushSender) *WorkerUsecase {
	return &WorkerUsecase{
		subRepo: subRepo,
		logRepo: logRepo,
		sender:  sender,
	}
}

func (u *WorkerUsecase) ProcessTask(ctx context.Context, task *kafka.PushJobTask) error {
	now := time.Now()

	log.Printf("[Worker] Processing push task for Subscription %s (Campaign: %s)", task.SubscriptionID, task.CampaignID)

	res := u.sender.SendPush(
		ctx,
		task.Endpoint,
		task.P256dhKey,
		task.AuthKey,
		task.Title,
		task.Body,
		task.IconURL,
		task.TargetURL,
		task.VAPIDPublicKey,
		task.VAPIDPrivKey,
		task.CampaignID.String(),
		task.SubscriptionID.String(),
	)

	status := domain.DeliveryStatusSent
	var errCode *int

	if res.StatusCode != 0 {
		code := res.StatusCode
		errCode = &code
		log.Printf("[Worker] Vendor Push Server Response HTTP %d for Endpoint: %s", res.StatusCode, task.Endpoint)
	}

	if res.Err != nil {
		status = domain.DeliveryStatusFailed
		log.Printf("[Worker] Push send error: %v", res.Err)
	}

	if res.IsExpired {
		log.Printf("[Worker] Subscription expired (HTTP %d). Soft-deleting sub: %s", res.StatusCode, task.SubscriptionID)
		_ = u.subRepo.Deactivate(ctx, task.SubscriptionID)
		status = domain.DeliveryStatusFailed
	}

	deliveryLog := &domain.DeliveryLog{
		ID:             uuid.New(),
		CampaignID:     task.CampaignID,
		SubscriptionID: task.SubscriptionID,
		Status:         status,
		ErrorCode:      errCode,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	return u.logRepo.BatchInsert(ctx, []*domain.DeliveryLog{deliveryLog})
}
