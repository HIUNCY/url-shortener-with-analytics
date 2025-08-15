package postgres

import (
	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type urlRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) domain.URLRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) Store(url *domain.URL) error {
	return r.db.Create(url).Error
}

func (r *urlRepository) FindByShortCode(shortCode string) (*domain.URL, error) {
	var url domain.URL
	err := r.db.Where("short_code = ?", shortCode).First(&url).Error
	return &url, err
}

func (r *urlRepository) FindByCustomAlias(customAlias string) (*domain.URL, error) {
	var url domain.URL
	err := r.db.Where("custom_alias = ?", customAlias).First(&url).Error
	return &url, err
}

// Implementasi method lain bisa ditambahkan nanti
func (r *urlRepository) FindAllByUserID(userID uuid.UUID) ([]domain.URL, error) { return nil, nil }
func (r *urlRepository) Update(url *domain.URL) error                           { return nil }
func (r *urlRepository) Delete(url *domain.URL) error                           { return nil }
