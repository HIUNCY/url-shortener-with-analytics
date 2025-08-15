package response

import (
	"time"
)

// APIKeyResponse adalah DTO untuk payload data pada respons regenerasi API Key.
type APIKeyResponse struct {
	APIKey      string    `json:"api_key"`
	GeneratedAt time.Time `json:"generated_at"`
}

// ProfileSuccessResponse adalah wrapper untuk Swagger.
type ProfileSuccessResponse struct {
	Success   bool         `json:"success" example:"true"`
	Data      UserResponse `json:"data"`
	Timestamp time.Time    `json:"timestamp"`
}

// APIKeySuccessResponse adalah wrapper untuk Swagger.
type APIKeySuccessResponse struct {
	Success   bool           `json:"success" example:"true"`
	Message   string         `json:"message" example:"API key regenerated successfully"`
	Data      APIKeyResponse `json:"data"`
	Timestamp time.Time      `json:"timestamp"`
}
