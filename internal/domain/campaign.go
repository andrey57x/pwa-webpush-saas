package domain

import (
	"time"

	"github.com/google/uuid"
)

type CampaignStatus string

const (
	CampaignStatusDraft      CampaignStatus = "DRAFT"
	CampaignStatusProcessing CampaignStatus = "PROCESSING"
	CampaignStatusCompleted  CampaignStatus = "COMPLETED"
	CampaignStatusFailed     CampaignStatus = "FAILED"
)

type Campaign struct {
	ID              uuid.UUID      `json:"id" db:"id"`
	AppID           uuid.UUID      `json:"app_id" db:"app_id"`
	CreatedByUserID *uuid.UUID     `json:"created_by_user_id,omitempty" db:"created_by_user_id"`
	Title           string         `json:"title" db:"title"`
	Body            string         `json:"body" db:"body"`
	IconURL         *string        `json:"icon_url,omitempty" db:"icon_url"`
	TargetURL       *string        `json:"target_url,omitempty" db:"target_url"`
	Status          CampaignStatus `json:"status" db:"status"`
	TotalTargeted   int            `json:"total_targeted" db:"total_targeted"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}
