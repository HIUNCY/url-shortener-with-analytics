package response

import (
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/google/uuid"
)

// CreateURLResponse adalah DTO untuk payload data pada respons pembuatan URL.
type CreateURLResponse struct {
	ID          uuid.UUID  `json:"id"`
	OriginalURL string     `json:"original_url"`
	ShortCode   string     `json:"short_code"`
	ShortURL    string     `json:"short_url"`
	CustomAlias *string    `json:"custom_alias,omitempty"`
	Title       *string    `json:"title,omitempty"`
	QRCode      string     `json:"qr_code"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// CreateURLSuccessResponse adalah wrapper untuk Swagger.
type CreateURLSuccessResponse struct {
	Success   bool              `json:"success" example:"true"`
	Message   string            `json:"message" example:"Short URL created successfully"`
	Data      CreateURLResponse `json:"data"`
	Timestamp time.Time         `json:"timestamp"`
}

func ToCreateURLResponse(url *domain.URL, shortURL, qrCode string) CreateURLResponse {
	return CreateURLResponse{
		ID:          url.ID,
		OriginalURL: url.OriginalURL,
		ShortCode:   url.ShortCode,
		ShortURL:    shortURL,
		CustomAlias: url.CustomAlias,
		Title:       url.Title,
		QRCode:      qrCode,
		ExpiresAt:   url.ExpiresAt,
		CreatedAt:   url.CreatedAt,
	}
}
