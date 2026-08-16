package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/andrey57x/pwa-webpush-saas/internal/delivery/middleware"
	"github.com/andrey57x/pwa-webpush-saas/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AppHandler struct {
	appUsecase AppUsecase
}

func NewAppHandler(appUsecase AppUsecase) *AppHandler {
	return &AppHandler{appUsecase: appUsecase}
}

func (h *AppHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantIDStr, _ := r.Context().Value(middleware.TenantIDKey).(string)
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil || tenantID == uuid.Nil {
		http.Error(w, `{"error":"tenant context missing"}`, http.StatusBadRequest)
		return
	}

	var req usecase.CreateAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json payload"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.AppCode == "" {
		http.Error(w, `{"error":"name and app_code are required"}`, http.StatusBadRequest)
		return
	}

	req.TenantID = tenantID

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

func (h *AppHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantIDStr, _ := r.Context().Value(middleware.TenantIDKey).(string)
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil || tenantID == uuid.Nil {
		http.Error(w, `{"error":"tenant context missing"}`, http.StatusBadRequest)
		return
	}

	apps, err := h.appUsecase.ListApps(r.Context(), tenantID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(apps)
}
