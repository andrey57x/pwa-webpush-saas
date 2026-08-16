package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/andrey57x/pwa-webpush-saas/internal/connectors/kafka"
	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/google/uuid"
)

type CampaignUsecase struct {
	campaignRepo CampaignRepository
	subRepo      SubscriptionRepository
	appRepo      AppRepository
	dedupRepo    DedupRepository
	producer     PushJobProducer
}

func NewCampaignUsecase(
	campaignRepo CampaignRepository,
	subRepo SubscriptionRepository,
	appRepo AppRepository,
	dedupRepo DedupRepository,
	producer PushJobProducer,
) *CampaignUsecase {
	return &CampaignUsecase{
		campaignRepo: campaignRepo,
		subRepo:      subRepo,
		appRepo:      appRepo,
		dedupRepo:    dedupRepo,
		producer:     producer,
	}
}

// LaunchCampaign выборку подписчиков пачками, фильтрует через Redis и отправляет в Kafka
func (u *CampaignUsecase) LaunchCampaign(ctx context.Context, campaignID uuid.UUID) error {
	campaign, err := u.campaignRepo.GetByID(ctx, campaignID)
	if err != nil || campaign == nil {
		return fmt.Errorf("campaign not found: %w", err)
	}

	app, err := u.appRepo.GetByID(ctx, campaign.AppID)
	if err != nil || app == nil {
		return fmt.Errorf("app not found: %w", err)
	}

	// Обновляем статус кампании на PROCESSING
	_ = u.campaignRepo.UpdateStatus(ctx, campaignID, domain.CampaignStatusProcessing, 0)

	batchSize := 500
	offset := 0
	totalSent := 0

	for {
		// 1. Извлекаем активные подписки батчами
		subs, err := u.subRepo.GetActiveByAppID(ctx, app.ID, batchSize, offset)
		if err != nil || len(subs) == 0 {
			break
		}

		subIDs := make([]uuid.UUID, len(subs))
		subMap := make(map[uuid.UUID]*domain.Subscription, len(subs))
		for i, sub := range subs {
			subIDs[i] = sub.ID
			subMap[sub.ID] = sub
		}

		// 2. Дедупликация через Redis Pipeline (SETNX c TTL 24 часа)
		uniqueIDs, err := u.dedupRepo.FilterDuplicates(ctx, campaignID, subIDs, 24*time.Hour)
		if err != nil {
			uniqueIDs = subIDs // В случае сбоя Redis продолжаем отправку
		}

		// 3. Формируем задачи для Kafka
		tasks := make([]*kafka.PushJobTask, 0, len(uniqueIDs))
		for _, id := range uniqueIDs {
			sub := subMap[id]
			icon, target := "", ""
			if campaign.IconURL != nil {
				icon = *campaign.IconURL
			}
			if campaign.TargetURL != nil {
				target = *campaign.TargetURL
			}

			tasks = append(tasks, &kafka.PushJobTask{
				CampaignID:     campaignID,
				SubscriptionID: sub.ID,
				Endpoint:       sub.Endpoint,
				P256dhKey:      sub.P256dhKey,
				AuthKey:        sub.AuthKey,
				Title:          campaign.Title,
				Body:           campaign.Body,
				IconURL:        icon,
				TargetURL:      target,
				VAPIDPublicKey: app.VAPIDPublicKey,
				VAPIDPrivKey:   app.VAPIDPrivateKey,
			})
		}

		// 4. Публикуем задачи в Kafka
		if len(tasks) > 0 {
			if err := u.producer.PublishPushJobs(ctx, tasks); err != nil {
				return fmt.Errorf("failed to publish push jobs: %w", err)
			}
			totalSent += len(tasks)
		}

		if len(subs) < batchSize {
			break
		}
		offset += batchSize
	}

	// Финализируем статус кампании
	return u.campaignRepo.UpdateStatus(ctx, campaignID, domain.CampaignStatusCompleted, totalSent)
}
