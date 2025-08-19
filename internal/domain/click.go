package domain

import (
	"time"

	"github.com/google/uuid"
)

type Click struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	URLID      uuid.UUID `gorm:"type:uuid;not null"`
	IPAddress  string
	UserAgent  string
	Referer    string
	Country    string
	Region     string
	City       string
	Browser    string
	OS         string
	DeviceType string
	IsUnique   bool `gorm:"default:false"`
	ClickedAt  time.Time
}

type TimeSeriesResult struct {
	Date  time.Time `gorm:"type:date"`
	Count int64
}

type GroupedResult struct {
	Value string
	Count int64
}

type ClickRepository interface {
	Store(click *Click) error
	GetTotalClicks(urlID uuid.UUID, since time.Time) (int64, error)
	GetTopReferrer(urlID uuid.UUID, since time.Time) (string, error)
	GetTopCountry(urlID uuid.UUID, since time.Time) (string, error)
	GetClicksOverTime(urlID uuid.UUID, since time.Time) ([]TimeSeriesResult, error)
	GetTopCountries(urlID uuid.UUID, since time.Time, limit int) ([]GroupedResult, error)
	GetTopReferrers(urlID uuid.UUID, since time.Time, limit int) ([]GroupedResult, error)
	GetDeviceStats(urlID uuid.UUID, since time.Time) ([]GroupedResult, error)
	GetBrowserStats(urlID uuid.UUID, since time.Time) ([]GroupedResult, error)
	GetOSStats(urlID uuid.UUID, since time.Time) ([]GroupedResult, error)
}
