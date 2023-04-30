package tokenutil

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/caarlos0/env"
	"github.com/google/uuid"
)

type Config struct {
	Key []byte `env:"PRIVATE_KEY"`
}

func GenerateUserId() string {
	id := uuid.New()
	return id.String()
}

func GenerateTokenForUser(userId string) string {
	var key []byte

	var cfg Config
	err := env.Parse(&cfg)
	if err == nil {
		key = cfg.Key
	}

	h := hmac.New(sha256.New, key)
	h.Write([]byte(userId))
	sign := h.Sum(nil)

	return userId + "." + hex.EncodeToString(sign)
}

func IsTokenValid(token string) bool {
	splitToken := strings.Split(token, ".")
	if len(splitToken) != 2 {
		return false
	}

	data, err := hex.DecodeString(splitToken[1])
	if err != nil {
		return false
	}

	userId := splitToken[0]

	h := hmac.New(sha256.New, nil)
	h.Write([]byte(userId))
	sign := h.Sum(nil)

	return hmac.Equal(sign, data)
}

func DecodeUserIdFromToken(token string) (string, error) {
	if token == "" {
		return "", errors.New("empty token")
	}
	splitToken := strings.Split(token, ".")
	if len(splitToken) != 2 {
		return "", errors.New("cant decode user id from token")
	}

	return splitToken[0], nil
}
