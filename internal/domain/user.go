package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleSuperAdmin  UserRole = "SUPER_ADMIN"
	RoleTenantAdmin UserRole = "TENANT_ADMIN"
	RoleTenantUser  UserRole = "TENANT_USER"
)

type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	TenantID     *uuid.UUID `json:"tenant_id,omitempty" db:"tenant_id"` // NULL для SUPER_ADMIN
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	Role         UserRole   `json:"role" db:"role"`
	Status       string     `json:"status" db:"status"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}
