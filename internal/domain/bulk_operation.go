package domain

import (
	"time"

	"github.com/google/uuid"
)

// BulkOperation merepresentasikan tugas pemrosesan URL secara massal.
type BulkOperation struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID         uuid.UUID `gorm:"type:uuid;not null"`
	OperationType  string    `gorm:"not null"`
	TotalCount     int       `gorm:"not null"`
	SuccessCount   int       `gorm:"default:0"`
	FailedCount    int       `gorm:"default:0"`
	Status         string    `gorm:"default:'pending'"`
	FilePath       *string
	ResultFilePath *string
	ErrorDetails   *string
	StartedAt      time.Time
	CompletedAt    *time.Time
}

// BulkOperationRepository mendefinisikan kontrak untuk interaksi data operasi massal.
type BulkOperationRepository interface {
	Store(op *BulkOperation) error
	FindByID(id uuid.UUID) (*BulkOperation, error)
	Update(op *BulkOperation) error
}
