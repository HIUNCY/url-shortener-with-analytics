package domain

import (
	"time"

	"github.com/google/uuid"
)

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
}

type FindAllOptions struct {
	Search string
	SortBy string
	Order  string
	Limit  int
	Offset int
}

type DashboardSummaryResult struct {
	TotalURLs   int64
	TotalClicks int64
	ActiveURLs  int64
}

type URLRepository interface {
	Store(url *URL) error
	FindByShortCode(shortCode string) (*URL, error)
	FindByCustomAlias(customAlias string) (*URL, error)
	FindByID(id uuid.UUID) (*URL, error)
	FindAllByUserID(userID uuid.UUID, options *FindAllOptions) ([]URL, int64, error)
	Update(url *URL) error
	Delete(url *URL) error
	IncrementClickCount(urlID uuid.UUID) error
	GetDashboardSummary(userID uuid.UUID) (*DashboardSummaryResult, error)
	GetTopPerformingURLs(userID uuid.UUID, limit int) ([]URL, error)
	GetRecentActivity(userID uuid.UUID, limit int) ([]URL, error)
}
