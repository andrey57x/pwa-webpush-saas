package usecase

import (
	"context"
	"errors"
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

func (u *CampaignUsecase) LaunchCampaign(ctx context.Context, tenantID, campaignID uuid.UUID) error {
	campaign, err := u.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("failed to query campaign: %w", err)
	}
	if campaign == nil {
		return errors.New("campaign not found")
	}

	app, err := u.appRepo.GetByID(ctx, campaign.AppID)
	if err != nil {
		return fmt.Errorf("failed to query app: %w", err)
	}
	if app == nil || app.TenantID != tenantID {
		return errors.New("access denied: campaign does not belong to this tenant")
	}

	_ = u.campaignRepo.UpdateStatus(ctx, campaignID, domain.CampaignStatusProcessing, 0)

	batchSize := 500
	offset := 0
	totalSent := 0

	for {
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

		uniqueIDs, err := u.dedupRepo.FilterDuplicates(ctx, campaignID, subIDs, 24*time.Hour)
		if err != nil {
			uniqueIDs = subIDs
		}

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

	return u.campaignRepo.UpdateStatus(ctx, campaignID, domain.CampaignStatusCompleted, totalSent)
}

func (u *CampaignUsecase) ListCampaigns(ctx context.Context, tenantID, appID uuid.UUID) ([]*domain.Campaign, error) {
	app, err := u.appRepo.GetByID(ctx, appID)
	if err != nil || app == nil || app.TenantID != tenantID {
		return nil, errors.New("access denied or app not found")
	}
	return u.campaignRepo.ListByAppID(ctx, appID)
}
