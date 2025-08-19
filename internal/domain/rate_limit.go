package domain

import (
	"time"

	"github.com/google/uuid"
)

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

type RateLimitRepository interface {
	Store(rateLimit *RateLimit) error
}
