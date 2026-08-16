package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/andrey57x/pwa-webpush-saas/internal/pkg/hash"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthUsecase struct {
	userRepo   UserRepository
	tenantRepo TenantRepository
	jwtSecret  string
}

func NewAuthUsecase(userRepo UserRepository, tenantRepo TenantRepository, jwtSecret string) *AuthUsecase {
	return &AuthUsecase{
		userRepo:   userRepo,
		tenantRepo: tenantRepo,
		jwtSecret:  jwtSecret,
	}
}

type RegisterRequest struct {
	TenantID *uuid.UUID      `json:"tenant_id,omitempty"`
	Email    string          `json:"email"`
	Password string          `json:"password"`
	Role     domain.UserRole `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	User        *domain.User `json:"user"`
}

func (u *AuthUsecase) Register(ctx context.Context, req *RegisterRequest) (*domain.User, error) {
	existing, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("user with email [%s] already exists", req.Email)
	}

	passwordHash, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var tenantID uuid.UUID

	// Если tenant_id не передан, атомарно создаем компанию для админа
	if req.TenantID == nil && req.Role == domain.RoleTenantAdmin {
		tenant := &domain.Tenant{
			ID:        uuid.New(),
			Name:      "Компания " + req.Email,
			Status:    domain.TenantStatusActive,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := u.tenantRepo.Create(ctx, tenant); err != nil {
			return nil, fmt.Errorf("failed to create tenant record: %w", err)
		}
		tenantID = tenant.ID
	} else if req.TenantID != nil {
		tenantID = *req.TenantID
	}

	user := &domain.User{
		ID:           uuid.New(),
		TenantID:     &tenantID,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         req.Role,
		Status:       "ACTIVE",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *AuthUsecase) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	if !hash.CheckPassword(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid email or password")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    string(user.Role),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	if user.TenantID != nil {
		claims["tenant_id"] = user.TenantID.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign jwt: %w", err)
	}

	return &AuthResponse{
		AccessToken: tokenString,
		User:        user,
	}, nil
}
