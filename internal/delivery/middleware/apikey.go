package middleware

import (
	"context"
	"net/http"

	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/andrey57x/pwa-webpush-saas/internal/pkg/hash"
)

type apiKeyRepo interface {
	GetByHash(ctx context.Context, keyHash string) (*domain.APIKey, error)
}

// APIKeyMiddleware проверяет заголовок X-API-Key для внешних серверов тенантов
func APIKeyMiddleware(repo apiKeyRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawKey := r.Header.Get("X-API-Key")
			if rawKey == "" {
				http.Error(w, `{"error":"unauthorized: missing X-API-Key header"}`, http.StatusUnauthorized)
				return
			}

			keyHash := hash.HashAPIKey(rawKey)
			apiKey, err := repo.GetByHash(r.Context(), keyHash)
			if err != nil || apiKey == nil {
				http.Error(w, `{"error":"unauthorized: invalid X-API-Key"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), TenantIDKey, apiKey.TenantID.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
