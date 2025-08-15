package domain

import (
	"time"

	"github.com/google/uuid"
)

// User merepresentasikan entitas pengguna dalam sistem.
type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email        string     `gorm:"unique;not null"`
	PasswordHash string     `gorm:"not null"`
	APIKey       string     `gorm:"unique;not null"`
	FirstName    *string
	LastName     *string
	IsActive     bool       `gorm:"default:true"`
	PlanType     string     `gorm:"default:'free'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLoginAt  *time.Time
}

// UserRepository mendefinisikan kontrak untuk interaksi data pengguna.
type UserRepository interface {
	Store(user *User) error
	FindByID(id uuid.UUID) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByAPIKey(apiKey string) (*User, error)
	Update(user *User) error
}
