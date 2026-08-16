package router

import (
	"net/http"

	"github.com/andrey57x/pwa-webpush-saas/internal/delivery"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(
	subHandler *delivery.SubscriptionHandler,
	campaignHandler *delivery.CampaignHandler,
	appHandler *delivery.AppHandler,
	feedbackHandler *delivery.FeedbackHandler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-API-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Отдача PWA фронтенда и Service Worker
	r.Handle("/sw.js", http.FileServer(http.Dir("./frontend/public")))
	r.Handle("/", http.FileServer(http.Dir("./frontend")))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/subscribe", subHandler.Subscribe)

		// Получение и создание приложений
		r.Post("/apps", appHandler.Create)
		r.Get("/apps/{code}", appHandler.GetByCode)

		// Feedback loop от Service Worker
		r.Post("/feedback/ping", feedbackHandler.Ping)

		// Кампании
		r.Route("/campaigns", func(r chi.Router) {
			r.Post("/", campaignHandler.Create)
			r.Post("/{id}/send", campaignHandler.Send)
		})
	})

	return r
}
