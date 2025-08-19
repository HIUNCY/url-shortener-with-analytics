package response

import (
	"time"
)

type APIKeyResponse struct {
	APIKey      string    `json:"api_key"`
	GeneratedAt time.Time `json:"generated_at"`
}

type ProfileSuccessResponse struct {
	Success   bool         `json:"success" example:"true"`
	Data      UserResponse `json:"data"`
	Timestamp time.Time    `json:"timestamp"`
}

type APIKeySuccessResponse struct {
	Success   bool           `json:"success" example:"true"`
	Message   string         `json:"message" example:"API key regenerated successfully"`
	Data      APIKeyResponse `json:"data"`
	Timestamp time.Time      `json:"timestamp"`
}
