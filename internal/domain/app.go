package domain

import (
	"time"

	"github.com/google/uuid"
)

type App struct {
	ID              uuid.UUID `json:"id" db:"id"`
	TenantID        uuid.UUID `json:"tenant_id" db:"tenant_id"`
	Name            string    `json:"name" db:"name"`
	AppCode         string    `json:"app_code" db:"app_code"` // Публичный кодовый идентификатор PWA
	VAPIDPublicKey  string    `json:"vapid_public_key" db:"vapid_public_key"`
	VAPIDPrivateKey string    `json:"-" db:"vapid_private_key"` // Зашифрован мастер-ключом
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
