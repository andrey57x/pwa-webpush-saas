package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/andrey57x/pwa-webpush-saas/internal/pkg/hash"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	TenantIDKey contextKey = "tenant_id"
	RoleKey     contextKey = "role"
)

type apiKeyRepository interface {
	GetByHash(ctx context.Context, keyHash string) (*domain.APIKey, error)
}

// UnifiedAuthMiddleware автоматически определяет JWT или X-API-Key
func UnifiedAuthMiddleware(jwtSecret string, apiKeyRepo apiKeyRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// 1. Проверка X-API-Key в заголовке
			if rawAPIKey := r.Header.Get("X-API-Key"); rawAPIKey != "" {
				keyHash := hash.HashAPIKey(rawAPIKey)
				apiKey, err := apiKeyRepo.GetByHash(ctx, keyHash)
				if err != nil || apiKey == nil {
					http.Error(w, `{"error":"unauthorized: invalid X-API-Key"}`, http.StatusUnauthorized)
					return
				}
				ctx = context.WithValue(ctx, TenantIDKey, apiKey.TenantID.String())
				ctx = context.WithValue(ctx, RoleKey, string(domain.RoleTenantAdmin))
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// 2. Проверка JWT Bearer Token
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString := strings.TrimPrefix(authHeader, "Bearer ")
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, http.ErrAbortHandler
					}
					return []byte(jwtSecret), nil
				})

				if err == nil && token.Valid {
					if claims, ok := token.Claims.(jwt.MapClaims); ok {
						if userID, exists := claims["user_id"].(string); exists {
							ctx = context.WithValue(ctx, UserIDKey, userID)
						}
						if tenantID, exists := claims["tenant_id"].(string); exists {
							ctx = context.WithValue(ctx, TenantIDKey, tenantID)
						}
						if role, exists := claims["role"].(string); exists {
							ctx = context.WithValue(ctx, RoleKey, role)
						}
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}

			http.Error(w, `{"error":"unauthorized: missing or invalid credentials"}`, http.StatusUnauthorized)
		})
	}
}
