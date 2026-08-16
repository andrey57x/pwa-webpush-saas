package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/andrey57x/pwa-webpush-saas/internal/usecase"
)

type FeedbackHandler struct {
	usecase *usecase.FeedbackUsecase
}

func NewFeedbackHandler(usecase *usecase.FeedbackUsecase) *FeedbackHandler {
	return &FeedbackHandler{usecase: usecase}
}

func (h *FeedbackHandler) Ping(w http.ResponseWriter, r *http.Request) {
	var req usecase.FeedbackPingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json payload"}`, http.StatusBadRequest)
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
