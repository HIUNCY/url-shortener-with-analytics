package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomString membuat string acak yang aman dengan panjang tertentu.
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateAPIKey membuat API Key baru untuk pengguna.
func GenerateAPIKey() (string, error) {
	// Menghasilkan 32-byte random string, aman untuk API Key
	return GenerateRandomString(32)
}

// GenerateShortCode membuat short code acak untuk URL.
func GenerateShortCode() (string, error) {
	// Menghasilkan 6-byte random string, cukup unik untuk short code
	// dan menghasilkan string sekitar 8 karakter
	return GenerateRandomString(6)
}
