package ece

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"
)

// EncryptPayload шифрует сообщение строго по стандартам RFC 8291 / RFC 8188 (aes128gcm)
func EncryptPayload(message []byte, p256dhBase64, authBase64 string) ([]byte, error) {
	// 1. Декодирование ключей клиента из Base64URL
	userPublicKeyBytes, err := base64.RawURLEncoding.DecodeString(p256dhBase64)
	if err != nil {
		return nil, fmt.Errorf("invalid user p256dh key: %w", err)
	}

	authSecretBytes, err := base64.RawURLEncoding.DecodeString(authBase64)
	if err != nil {
		return nil, fmt.Errorf("invalid user auth key: %w", err)
	}

	// 2. Генерация 16 байт случайной соли (Salt)
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// 3. Создание эфеемерного ключа ECDH P-256
	curve := ecdh.P256()
	localPrivateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ecdh key: %w", err)
	}
	localPublicKey := localPrivateKey.PublicKey()

	remotePublicKey, err := curve.NewPublicKey(userPublicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse remote public key: %w", err)
	}

	// 4. ECDH Shared Secret
	sharedSecret, err := localPrivateKey.ECDH(remotePublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to compute ecdh shared secret: %w", err)
	}

	// 5. Раздельная деривация Extract и Expand по RFC 8291
	// PRK_key = HKDF-Extract(salt = authSecretBytes, secret = sharedSecret)
	prkKey := hkdf.Extract(sha256.New, sharedSecret, authSecretBytes)

	// IKM = HKDF-Expand(PRK = prkKey, info = "WebPush: info\x00" || remote_pk || local_pk, L = 32)
	ikmInfo := append([]byte("WebPush: info\x00"), remotePublicKey.Bytes()...)
	ikmInfo = append(ikmInfo, localPublicKey.Bytes()...)

	ikmReader := hkdf.Expand(sha256.New, prkKey, ikmInfo)
	ikm := make([]byte, 32)
	if _, err := io.ReadFull(ikmReader, ikm); err != nil {
		return nil, fmt.Errorf("failed to derive ikm: %w", err)
	}

	// PRK_ece = HKDF-Extract(salt = salt, secret = ikm) по RFC 8188 Section 2.2
	prkEce := hkdf.Extract(sha256.New, ikm, salt)

	// CEK = HKDF-Expand(PRK_ece, info = "Content-Encoding: aes128gcm\x00", L = 16)
	cekReader := hkdf.Expand(sha256.New, prkEce, []byte("Content-Encoding: aes128gcm\x00"))
	cek := make([]byte, 16)
	if _, err := io.ReadFull(cekReader, cek); err != nil {
		return nil, fmt.Errorf("failed to derive cek: %w", err)
	}

	// Nonce = HKDF-Expand(PRK_ece, info = "Content-Encoding: nonce\x00", L = 12)
	nonceReader := hkdf.Expand(sha256.New, prkEce, []byte("Content-Encoding: nonce\x00"))
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(nonceReader, nonce); err != nil {
		return nil, fmt.Errorf("failed to derive nonce: %w", err)
	}

	// 6. Record Padding по RFC 8188 (добавляем 0x02 в конец единственной записи)
	paddedPayload := append(message, 0x02)

	// 7. Шифрование AES-128-GCM
	block, err := aes.NewCipher(cek)
	if err != nil {
		return nil, fmt.Errorf("failed to create aes cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create gcm: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, paddedPayload, nil)

	// 8. Сборка RFC 8188 Header (Salt[16] + RS[4] + IDLen[1] + EphemeralPubKey[65] + Ciphertext)
	recordSize := make([]byte, 4)
	binary.BigEndian.PutUint32(recordSize, 4096)

	header := make([]byte, 0, 16+4+1+65+len(ciphertext))
	header = append(header, salt...)
	header = append(header, recordSize...)
	header = append(header, byte(len(localPublicKey.Bytes())))
	header = append(header, localPublicKey.Bytes()...)
	header = append(header, ciphertext...)

	return header, nil
}
