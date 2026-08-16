package pushsender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/andrey57x/pwa-webpush-saas/internal/pkg/crypto/ece"
	"github.com/andrey57x/pwa-webpush-saas/internal/pkg/crypto/vapid"
)

type PushResult struct {
	StatusCode  int
	IsExpired   bool
	IsRetryable bool
	Err         error
}

type Sender struct {
	httpClient   *http.Client
	contactEmail string
}

func NewSender(contactEmail string) *Sender {
	return &Sender{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		contactEmail: contactEmail,
	}
}

func (s *Sender) SendPush(
	ctx context.Context,
	endpoint, p256dh, auth, title, body, iconURL, targetURL, vapidPublicKey, vapidPrivateKey, campaignID, subscriptionID string,
) *PushResult {
	payloadMap := map[string]string{
		"title":           title,
		"body":            body,
		"icon":            iconURL,
		"url":             targetURL,
		"campaign_id":     campaignID,
		"subscription_id": subscriptionID,
	}
	payloadJSON, err := json.Marshal(payloadMap)
	if err != nil {
		return &PushResult{StatusCode: 0, Err: fmt.Errorf("failed to marshal payload: %w", err)}
	}

	encryptedPayload, err := ece.EncryptPayload(payloadJSON, p256dh, auth)
	if err != nil {
		return &PushResult{StatusCode: 0, Err: fmt.Errorf("ece encryption failed: %w", err)}
	}

	authHeader, err := vapid.CreateVAPIDHeader(endpoint, s.contactEmail, vapidPublicKey, vapidPrivateKey)
	if err != nil {
		return &PushResult{StatusCode: 0, Err: fmt.Errorf("vapid header failed: %w", err)}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(encryptedPayload))
	if err != nil {
		return &PushResult{StatusCode: 0, Err: err}
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Encoding", "aes128gcm")
	req.Header.Set("TTL", "86400")
	req.Header.Set("Urgency", "high")
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Crypto-Key", "p256ecdsa="+strings.TrimRight(vapidPublicKey, "="))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return &PushResult{StatusCode: 0, IsRetryable: true, Err: err}
	}
	defer resp.Body.Close()

	result := &PushResult{StatusCode: resp.StatusCode}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		return result
	case http.StatusGone, http.StatusNotFound:
		result.IsExpired = true
		return result
	case http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		result.IsRetryable = true
		return result
	default:
		result.Err = fmt.Errorf("push server returned status: %d", resp.StatusCode)
		return result
	}
}
