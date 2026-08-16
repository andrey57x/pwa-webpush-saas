package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/andrey57x/pwa-webpush-saas/internal/usecase"
	"github.com/google/uuid"
)

type FeedbackHandler struct {
	usecase FeedbackUsecase
}

func NewFeedbackHandler(usecase FeedbackUsecase) *FeedbackHandler {
	return &FeedbackHandler{usecase: usecase}
}

func (h *FeedbackHandler) Ping(w http.ResponseWriter, r *http.Request) {
	var req usecase.FeedbackPingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json payload"}`, http.StatusBadRequest)
		return
	}

	if req.CampaignID == uuid.Nil || req.SubscriptionID == uuid.Nil {
		http.Error(w, `{"error":"campaign_id and subscription_id are required"}`, http.StatusBadRequest)
		return
	}

	if err := h.usecase.ProcessFeedback(r.Context(), &req); err != nil {
		http.Error(w, `{"error":"failed to record feedback"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"recorded"}`))
}
