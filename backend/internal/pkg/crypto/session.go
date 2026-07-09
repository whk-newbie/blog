package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/whk-newbie/blog/internal/pkg/redis"
)

const (
	SessionKeyPrefix = "session:enc:"
	SessionTTL       = 30 * time.Minute
	AESKeySize       = 32
)

// GenerateSessionID generates a new session ID
func GenerateSessionID() string {
	return uuid.New().String()
}

// GenerateAESKey generates a random AES-256 key
func GenerateAESKey() ([]byte, error) {
	key := make([]byte, AESKeySize)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %w", err)
	}
	return key, nil
}

// StoreSessionKey stores session -> AES key in Redis
func StoreSessionKey(sessionID string, aesKey []byte) error {
	key := SessionKeyPrefix + sessionID
	return redis.Set(key, base64.StdEncoding.EncodeToString(aesKey), SessionTTL)
}

// GetSessionKey retrieves AES key from Redis
func GetSessionKey(sessionID string) ([]byte, error) {
	key := SessionKeyPrefix + sessionID
	val, err := redis.GetValue(key)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(val)
}

// RefreshSessionKey refreshes session TTL
func RefreshSessionKey(sessionID string) error {
	key := SessionKeyPrefix + sessionID
	return redis.Expire(key, SessionTTL)
}
