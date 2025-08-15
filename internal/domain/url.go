package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// URL adalah entitas inti yang merepresentasikan URL yang dipersingkat.
type URL struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID           *uuid.UUID `gorm:"type:uuid"`
	OriginalURL      string     `gorm:"not null"`
	ShortCode        string     `gorm:"unique;not null"`
	CustomAlias      *string    `gorm:"unique"`
	DomainID         *uuid.UUID `gorm:"type:uuid"`
	Title            *string
	Description      *string
	PasswordHash     *string
	IsActive         bool `gorm:"default:true"`
	ClickCount       int  `gorm:"default:0"`
	UniqueClickCount int  `gorm:"default:0"`
	ExpiresAt        *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	LastClickedAt    *time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

// URLRepository mendefinisikan kontrak untuk interaksi data URL.
type URLRepository interface {
	Store(url *URL) error
	FindByShortCode(shortCode string) (*URL, error)
	FindByCustomAlias(customAlias string) (*URL, error)
	FindAllByUserID(userID uuid.UUID) ([]URL, error)
	Update(url *URL) error
	Delete(url *URL) error
}
