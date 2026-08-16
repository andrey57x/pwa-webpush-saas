package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/google/uuid"
)

type AnalyticsUsecase struct {
	analyticsRepo AnalyticsRepository
	campaignRepo  CampaignRepository
	appRepo       AppRepository
}

func NewAnalyticsUsecase(analyticsRepo AnalyticsRepository, campaignRepo CampaignRepository, appRepo AppRepository) *AnalyticsUsecase {
	return &AnalyticsUsecase{
		analyticsRepo: analyticsRepo,
		campaignRepo:  campaignRepo,
		appRepo:       appRepo,
	}
}

func (u *AnalyticsUsecase) GetCampaignStats(ctx context.Context, tenantID, campaignID uuid.UUID) (*domain.CampaignStats, error) {
	campaign, err := u.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to query campaign: %w", err)
	}
	if campaign == nil {
		return nil, errors.New("campaign not found")
	}

	app, err := u.appRepo.GetByID(ctx, campaign.AppID)
	if err != nil || app == nil || app.TenantID != tenantID {
		return nil, errors.New("access denied: campaign does not belong to this tenant")
	}

	counts, err := u.analyticsRepo.GetCampaignStatusCounts(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate status counts: %w", err)
	}

	sent := counts[domain.DeliveryStatusSent]
	delivered := counts[domain.DeliveryStatusDelivered]
	clicked := counts[domain.DeliveryStatusClicked]
	failed := counts[domain.DeliveryStatusFailed]

	var ctr float64
	if delivered > 0 {
		ctr = (float64(clicked) / float64(delivered)) * 100.0
	} else if sent > 0 {
		ctr = (float64(clicked) / float64(sent)) * 100.0
	}
	ctr = math.Round(ctr*100) / 100

	return &domain.CampaignStats{
		CampaignID:     campaign.ID,
		Title:          campaign.Title,
		Status:         string(campaign.Status),
		TotalTargeted:  campaign.TotalTargeted,
		SentCount:      sent,
		DeliveredCount: delivered,
		ClickedCount:   clicked,
		FailedCount:    failed,
		CTR:            ctr,
	}, nil
}
