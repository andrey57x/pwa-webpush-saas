package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/andrey57x/pwa-webpush-saas/internal/delivery/middleware"
	"github.com/google/uuid"
)

type APIKeyHandler struct {
	apiKeyUsecase APIKeyUsecase // Интерфейс из interfaces.go!
}

func NewAPIKeyHandler(apiKeyUsecase APIKeyUsecase) *APIKeyHandler {
	return &APIKeyHandler{apiKeyUsecase: apiKeyUsecase}
}

type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}

func (h *APIKeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantIDStr, _ := r.Context().Value(middleware.TenantIDKey).(string)
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil || tenantID == uuid.Nil {
		http.Error(w, `{"error":"tenant context missing"}`, http.StatusBadRequest)
		return
	}

	var req CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, `{"error":"name is required"}`, http.StatusBadRequest)
		return
	}

	res, err := h.apiKeyUsecase.CreateKey(r.Context(), tenantID, req.Name)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
}

func (h *APIKeyHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantIDStr, _ := r.Context().Value(middleware.TenantIDKey).(string)
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil || tenantID == uuid.Nil {
		http.Error(w, `{"error":"tenant context missing"}`, http.StatusBadRequest)
		return
	}

	keys, err := h.apiKeyUsecase.ListKeys(r.Context(), tenantID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(keys)
}
