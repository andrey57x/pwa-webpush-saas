package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	AppID          uuid.UUID  `json:"app_id" db:"app_id"`
	UserIdentifier *string    `json:"user_identifier,omitempty" db:"user_identifier"`
	Endpoint       string     `json:"endpoint" db:"endpoint"`
	P256dhKey      string     `json:"p256dh_key" db:"p256dh_key"`
	AuthKey        string     `json:"auth_key" db:"auth_key"`
	Browser        string     `json:"browser" db:"browser"`
	OS             string     `json:"os" db:"os"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	DeactivatedAt  *time.Time `json:"deactivated_at,omitempty" db:"deactivated_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}
