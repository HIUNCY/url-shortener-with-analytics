package response

import (
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/google/uuid"
)

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

type CreateURLSuccessResponse struct {
	Success   bool              `json:"success" example:"true"`
	Message   string            `json:"message" example:"Short URL created successfully"`
	Data      CreateURLResponse `json:"data"`
	Timestamp time.Time         `json:"timestamp"`
}

type URLListItemResponse struct {
	ID          uuid.UUID  `json:"id"`
	OriginalURL string     `json:"original_url"`
	ShortCode   string     `json:"short_code"`
	ShortURL    string     `json:"short_url"`
	Title       *string    `json:"title,omitempty"`
	ClickCount  int        `json:"click_count"`
	IsActive    bool       `json:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type URLListResponse struct {
	URLs       []URLListItemResponse `json:"urls"`
	Pagination PaginationResponse    `json:"pagination"`
}

type URLListSuccessResponse struct {
	Success   bool            `json:"success" example:"true"`
	Data      URLListResponse `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
}

type URLDetailsResponse struct {
	ID                  uuid.UUID  `json:"id"`
	OriginalURL         string     `json:"original_url"`
	ShortCode           string     `json:"short_code"`
	ShortURL            string     `json:"short_url"`
	CustomAlias         *string    `json:"custom_alias,omitempty"`
	Title               *string    `json:"title,omitempty"`
	Description         *string    `json:"description,omitempty"`
	ClickCount          int        `json:"click_count"`
	UniqueClickCount    int        `json:"unique_click_count"`
	IsActive            bool       `json:"is_active"`
	IsPasswordProtected bool       `json:"is_password_protected"`
	ExpiresAt           *time.Time `json:"expires_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	LastClickedAt       *time.Time `json:"last_clicked_at,omitempty"`
}

type URLDetailsSuccessResponse struct {
	Success   bool               `json:"success" example:"true"`
	Data      URLDetailsResponse `json:"data"`
	Timestamp time.Time          `json:"timestamp"`
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

func ToURLDetailsResponse(url *domain.URL, shortURL string) URLDetailsResponse {
	return URLDetailsResponse{
		ID:                  url.ID,
		OriginalURL:         url.OriginalURL,
		ShortCode:           url.ShortCode,
		ShortURL:            shortURL,
		CustomAlias:         url.CustomAlias,
		Title:               url.Title,
		Description:         url.Description,
		ClickCount:          url.ClickCount,
		UniqueClickCount:    url.UniqueClickCount,
		IsActive:            url.IsActive,
		IsPasswordProtected: url.PasswordHash != nil,
		ExpiresAt:           url.ExpiresAt,
		CreatedAt:           url.CreatedAt,
		UpdatedAt:           url.UpdatedAt,
		LastClickedAt:       url.LastClickedAt,
	}
}
