package postgres

import (
	"fmt"
	"strings"

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

func (r *urlRepository) FindByID(id uuid.UUID) (*domain.URL, error) {
	var url domain.URL
	err := r.db.Where("id = ?", id).First(&url).Error
	return &url, err
}

func (r *urlRepository) FindAllByUserID(userID uuid.UUID, options *domain.FindAllOptions) ([]domain.URL, int64, error) {
	var urls []domain.URL
	var total int64

	// Mulai query dasar
	query := r.db.Model(&domain.URL{}).Where("user_id = ?", userID)

	// Terapkan filter pencarian
	if options.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(options.Search))
		query = query.Where("LOWER(title) LIKE ? OR LOWER(original_url) LIKE ?", searchQuery, searchQuery)
	}

	// Hitung total data sebelum paginasi
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Terapkan sorting
	if options.SortBy != "" && options.Order != "" {
		query = query.Order(fmt.Sprintf("%s %s", options.SortBy, options.Order))
	} else {
		query = query.Order("created_at desc") // Default sort
	}

	// Terapkan paginasi (limit dan offset)
	query = query.Limit(options.Limit).Offset(options.Offset)

	// Eksekusi query untuk mendapatkan data
	if err := query.Find(&urls).Error; err != nil {
		return nil, 0, err
	}

	return urls, total, nil
}

func (r *urlRepository) Update(url *domain.URL) error {
	return r.db.Save(url).Error
}

func (r *urlRepository) Delete(url *domain.URL) error {
	return r.db.Delete(url).Error
}
