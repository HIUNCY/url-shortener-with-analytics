package request

import "time"

// CreateURLRequest adalah DTO untuk request pembuatan URL pendek.
type CreateURLRequest struct {
	OriginalURL string     `json:"original_url" binding:"required,url"`
	CustomAlias *string    `json:"custom_alias,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Password    *string    `json:"password,omitempty"`
}

type UpdateURLRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	IsActive    *bool      `json:"is_active,omitempty"`
}
