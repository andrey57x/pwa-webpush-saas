package vapid

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"time"
)

// VAPIDKeys содержит парные VAPID-ключи в формате Base64URL
type VAPIDKeys struct {
	PublicKey  string
	PrivateKey string
}

// GenerateKeys генерирует VAPID-ключи (ECDSA P-256) с использованием современного crypto/ecdh (без deprecated функций)
func GenerateKeys() (*VAPIDKeys, error) {
	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ecdsa key: %w", err)
	}

	// Переходим в crypto/ecdh для константной по времени работы
	ecdhPriv, err := ecdsaKey.ECDH()
	if err != nil {
		return nil, fmt.Errorf("failed to convert ecdsa to ecdh: %w", err)
	}

	// Публичный ключ ANSI X9.62 uncompressed point (65 байт: 0x04 || X || Y)
	pubBytes := ecdhPriv.PublicKey().Bytes()
	pubBase64 := base64.RawURLEncoding.EncodeToString(pubBytes)

	// Приватный ключ (ровно 32 байта)
	privBytes := ecdhPriv.Bytes()
	privBase64 := base64.RawURLEncoding.EncodeToString(privBytes)

	return &VAPIDKeys{
		PublicKey:  pubBase64,
		PrivateKey: privBase64,
	}, nil
}

// CreateVAPIDHeader формирует заголовок Authorization по RFC 8292
func CreateVAPIDHeader(endpointURL, contactEmail, publicKeyBase64, privateKeyBase64 string) (string, error) {
	parsedURL, err := url.Parse(endpointURL)
	if err != nil {
		return "", fmt.Errorf("invalid endpoint url: %w", err)
	}

	origin := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	// 1. Создаем JWT Header
	jwtHeader, err := json.Marshal(map[string]string{
		"typ": "JWT",
		"alg": "ES256",
	})
	if err != nil {
		return "", err
	}

	// 2. Создаем JWT Payload (12 часов)
	jwtPayload, err := json.Marshal(map[string]interface{}{
		"aud": origin,
		"exp": time.Now().Add(12 * time.Hour).Unix(),
		"sub": "mailto:" + contactEmail,
	})
	if err != nil {
		return "", err
	}

	headerB64 := base64.RawURLEncoding.EncodeToString(jwtHeader)
	payloadB64 := base64.RawURLEncoding.EncodeToString(jwtPayload)
	unsignedToken := headerB64 + "." + payloadB64

	// 3. Декодируем приватный ключ из Base64URL
	privBytes, err := base64.RawURLEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return "", fmt.Errorf("invalid private key base64: %w", err)
	}

	// Безопасное восстановление координат точки через crypto/ecdh без ScalarBaseMult
	ecdhPriv, err := ecdh.P256().NewPrivateKey(privBytes)
	if err != nil {
		return "", fmt.Errorf("invalid ecdh private key bytes: %w", err)
	}
	pubBytes := ecdhPriv.PublicKey().Bytes()

	privKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     new(big.Int).SetBytes(pubBytes[1:33]),  // X координата (32B)
			Y:     new(big.Int).SetBytes(pubBytes[33:65]), // Y координата (32B)
		},
		D: new(big.Int).SetBytes(privBytes),
	}

	// 4. Подписываем SHA256 хэш unsignedToken
	hash := sha256.Sum256([]byte(unsignedToken))
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign vapid jwt: %w", err)
	}

	// IEEE P1363 подпись (64 байта)
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	sigBytes := make([]byte, 64)
	copy(sigBytes[32-len(rBytes):32], rBytes)
	copy(sigBytes[64-len(sBytes):64], sBytes)

	sigB64 := base64.RawURLEncoding.EncodeToString(sigBytes)
	jwtToken := unsignedToken + "." + sigB64

	return fmt.Sprintf("vapid t=%s, k=%s", jwtToken, strings.TrimRight(publicKeyBase64, "=")), nil
}
