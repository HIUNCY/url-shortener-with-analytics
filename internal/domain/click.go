package domain

import (
	"time"

	"github.com/google/uuid"
)

// Click merepresentasikan data analitik dari setiap klik pada sebuah URL.
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

// ClickRepository mendefinisikan kontrak untuk penyimpanan data klik.
type ClickRepository interface {
	Store(click *Click) error
}
