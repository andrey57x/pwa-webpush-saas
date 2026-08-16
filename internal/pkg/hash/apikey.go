package hash

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

// GenerateAPIKey генерирует сырой ключ pk_live_... для пользователя и его SHA-256 хэш для БД
func GenerateAPIKey() (rawKey string, keyHash string, err error) {
	bytes := make([]byte, 24)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	rawKey = "pk_live_" + hex.EncodeToString(bytes) + uuid.New().String()[:8]

	hash := sha256.Sum256([]byte(rawKey))
	keyHash = hex.EncodeToString(hash[:])
	return rawKey, keyHash, nil
}

func HashAPIKey(rawKey string) string {
	hash := sha256.Sum256([]byte(rawKey))
	return hex.EncodeToString(hash[:])
}
