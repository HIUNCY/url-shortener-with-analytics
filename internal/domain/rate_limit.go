package domain

import (
	"time"

	"github.com/google/uuid"
)

// RateLimit merepresentasikan catatan penggunaan API untuk rate limiting.
type RateLimit struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID       *uuid.UUID `gorm:"type:uuid"`
	APIKey       *string
	IPAddress    string
	Endpoint     string
	RequestCount int `gorm:"default:1"`
	WindowStart  time.Time
	CreatedAt    time.Time
}

// RateLimitRepository mendefinisikan kontrak untuk interaksi data rate limit.
type RateLimitRepository interface {
	Store(rateLimit *RateLimit) error
}
