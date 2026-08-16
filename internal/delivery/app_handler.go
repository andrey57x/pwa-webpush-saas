package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/andrey57x/pwa-webpush-saas/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type AppHandler struct {
	appUsecase *usecase.AppUsecase
}

func NewAppHandler(appUsecase *usecase.AppUsecase) *AppHandler {
	return &AppHandler{appUsecase: appUsecase}
}

func (h *AppHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req usecase.CreateAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json payload"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.AppCode == "" {
		http.Error(w, `{"error":"name and app_code are required"}`, http.StatusBadRequest)
		return
	}

	app, err := h.appUsecase.CreateApp(r.Context(), &req)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(app)
}

func (h *AppHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	appCode := chi.URLParam(r, "code")
	if appCode == "" {
		http.Error(w, `{"error":"app code is required"}`, http.StatusBadRequest)
		return
	}

	app, err := h.appUsecase.GetByCode(r.Context(), appCode)
	if err != nil || app == nil {
		http.Error(w, `{"error":"app not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"app_code":         app.AppCode,
		"name":             app.Name,
		"vapid_public_key": app.VAPIDPublicKey,
	})
}
