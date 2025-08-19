package domain

import (
	"time"

	"github.com/google/uuid"
)

type Domain struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID            uuid.UUID `gorm:"type:uuid;not null"`
	DomainName        string    `gorm:"unique;not null"`
	IsVerified        bool      `gorm:"default:false"`
	VerificationToken string
	IsActive          bool `gorm:"default:true"`
	CreatedAt         time.Time
	VerifiedAt        *time.Time
}

type DomainRepository interface {
	Store(domain *Domain) error
	FindByDomainName(name string) (*Domain, error)
	FindAllByUserID(userID uuid.UUID) ([]Domain, error)
}
