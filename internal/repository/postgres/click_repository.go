package postgres

import (
	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"gorm.io/gorm"
)

type clickRepository struct {
	db *gorm.DB
}

// NewClickRepository membuat instance baru dari clickRepository.
func NewClickRepository(db *gorm.DB) domain.ClickRepository {
	return &clickRepository{db: db}
}

// Store menyimpan data klik baru ke database.
func (r *clickRepository) Store(click *domain.Click) error {
	return r.db.Create(click).Error
}
