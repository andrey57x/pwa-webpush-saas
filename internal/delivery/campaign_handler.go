package delivery

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/andrey57x/pwa-webpush-saas/internal/domain"
	"github.com/andrey57x/pwa-webpush-saas/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CampaignHandler struct {
	usecase      *usecase.CampaignUsecase
	campaignRepo usecase.CampaignRepository
}

func NewCampaignHandler(usecase *usecase.CampaignUsecase, campaignRepo usecase.CampaignRepository) *CampaignHandler {
	return &CampaignHandler{
		usecase:      usecase,
		campaignRepo: campaignRepo,
	}
}

type CreateCampaignRequest struct {
	AppID     uuid.UUID `json:"app_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	IconURL   *string   `json:"icon_url,omitempty"`
	TargetURL *string   `json:"target_url,omitempty"`
}

func (h *CampaignHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json payload"}`, http.StatusBadRequest)
		return
	}

	campaign := &domain.Campaign{
		ID:        uuid.New(),
		AppID:     req.AppID,
		Title:     req.Title,
		Body:      req.Body,
		IconURL:   req.IconURL,
		TargetURL: req.TargetURL,
		Status:    domain.CampaignStatusDraft,
	}

	if err := h.campaignRepo.Create(r.Context(), campaign); err != nil {
		http.Error(w, `{"error":"failed to create campaign"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(campaign)
}

func (h *CampaignHandler) Send(w http.ResponseWriter, r *http.Request) {
	campaignIDStr := chi.URLParam(r, "id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		http.Error(w, `{"error":"invalid campaign id"}`, http.StatusBadRequest)
		return
	}

	// Передаем context.Background() для асинхронной рассылки
	go func() {
		if err := h.usecase.LaunchCampaign(context.Background(), campaignID); err != nil {
			log.Printf("[Campaign] Error launching campaign %s: %v", campaignID, err)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte(`{"message":"campaign sending started"}`))
}
