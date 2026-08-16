package domain

import (
	"time"

	"github.com/google/uuid"
)

type DeliveryStatus string

const (
	DeliveryStatusSent      DeliveryStatus = "SENT"
	DeliveryStatusDelivered DeliveryStatus = "DELIVERED"
	DeliveryStatusClicked   DeliveryStatus = "CLICKED"
	DeliveryStatusFailed    DeliveryStatus = "FAILED"
)

type DeliveryLog struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	CampaignID     uuid.UUID      `json:"campaign_id" db:"campaign_id"`
	SubscriptionID uuid.UUID      `json:"subscription_id" db:"subscription_id"`
	Status         DeliveryStatus `json:"status" db:"status"`
	ErrorCode      *int           `json:"error_code,omitempty" db:"error_code"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}
