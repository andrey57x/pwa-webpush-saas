package domain

import (
	"github.com/google/uuid"
)

type CampaignStats struct {
	CampaignID     uuid.UUID `json:"campaign_id"`
	Title          string    `json:"title"`
	Status         string    `json:"status"`
	TotalTargeted  int       `json:"total_targeted"`
	SentCount      int       `json:"sent_count"`
	DeliveredCount int       `json:"delivered_count"`
	ClickedCount   int       `json:"clicked_count"`
	FailedCount    int       `json:"failed_count"`
	CTR            float64   `json:"ctr"` // Click-Through Rate в %
}
