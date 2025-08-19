package domain

import (
	"time"

	"github.com/google/uuid"
)

type QRCode struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	URLID     uuid.UUID `gorm:"type:uuid;not null"`
	QRData    string    `gorm:"not null"`
	Format    string    `gorm:"default:'png'"`
	Size      int       `gorm:"default:200"`
	CreatedAt time.Time
}

type QRCodeRepository interface {
	Store(qrCode *QRCode) error
	FindByURLID(urlID uuid.UUID) (*QRCode, error)
}
