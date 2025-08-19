package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func GenerateAPIKey() (string, error) {
	return GenerateRandomString(32)
}

func GenerateShortCode() (string, error) {
	return GenerateRandomString(6)
}
