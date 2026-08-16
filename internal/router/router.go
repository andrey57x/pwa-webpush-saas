package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/andrey57x/pwa-webpush-saas/internal/delivery"
	"github.com/andrey57x/pwa-webpush-saas/internal/delivery/middleware"
	"github.com/andrey57x/pwa-webpush-saas/internal/usecase"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func serveSPA(staticDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Clean(r.URL.Path)

		if strings.HasPrefix(filepath.Base(path), ".") {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		fullPath := filepath.Join(staticDir, path)
		info, err := os.Stat(fullPath)
		if err == nil && !info.IsDir() {
			http.ServeFile(w, r, fullPath)
			return
		}

		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	}
}

func NewRouter(
	subHandler *delivery.SubscriptionHandler,
	campaignHandler *delivery.CampaignHandler,
	appHandler *delivery.AppHandler,
	feedbackHandler *delivery.FeedbackHandler,
	authHandler *delivery.AuthHandler,
	analyticsHandler *delivery.AnalyticsHandler,
	apiKeyHandler *delivery.APIKeyHandler,
	apiKeyRepo usecase.APIKeyRepository,
	jwtSecret string,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-API-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/subscribe", subHandler.Subscribe)
		r.Get("/apps/{code}", appHandler.GetByCode)
		r.Post("/feedback/ping", feedbackHandler.Ping)

		r.Group(func(r chi.Router) {
			r.Use(middleware.UnifiedAuthMiddleware(jwtSecret, apiKeyRepo))

			r.Get("/apps", appHandler.List)
			r.Post("/apps", appHandler.Create)
			r.Post("/api-keys", apiKeyHandler.Create)
			r.Get("/api-keys", apiKeyHandler.List)

			r.Route("/campaigns", func(r chi.Router) {
				r.Get("/", campaignHandler.List)
				r.Post("/", campaignHandler.Create)
				r.Post("/{id}/send", campaignHandler.Send)
			})

			r.Get("/analytics/campaigns/{id}", analyticsHandler.GetCampaignStats)
		})
	})

	r.NotFound(serveSPA("./frontend"))

	return r
}
