package postgres

import (
	"fmt"
	"strings"
	"time"

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

	query := r.db.Model(&domain.URL{}).Where("user_id = ?", userID)

	if options.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(options.Search))
		query = query.Where("LOWER(title) LIKE ? OR LOWER(original_url) LIKE ?", searchQuery, searchQuery)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if options.SortBy != "" && options.Order != "" {
		query = query.Order(fmt.Sprintf("%s %s", options.SortBy, options.Order))
	} else {
		query = query.Order("created_at desc")
	}

	query = query.Limit(options.Limit).Offset(options.Offset)

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

func (r *urlRepository) IncrementClickCount(urlID uuid.UUID) error {
	return r.db.Model(&domain.URL{}).Where("id = ?", urlID).Updates(map[string]interface{}{
		"click_count":     gorm.Expr("click_count + 1"),
		"last_clicked_at": time.Now(),
	}).Error
}

func (r *urlRepository) GetDashboardSummary(userID uuid.UUID) (*domain.DashboardSummaryResult, error) {
	var result domain.DashboardSummaryResult
	err := r.db.Model(&domain.URL{}).
		Select("COUNT(*) as total_urls, COALESCE(SUM(click_count), 0) as total_clicks, COUNT(CASE WHEN is_active = true AND (expires_at IS NULL OR expires_at > NOW()) THEN 1 END) as active_urls").
		Where("user_id = ?", userID).
		Scan(&result).Error
	return &result, err
}

func (r *urlRepository) GetTopPerformingURLs(userID uuid.UUID, limit int) ([]domain.URL, error) {
	var urls []domain.URL
	err := r.db.Where("user_id = ?", userID).
		Order("click_count DESC").
		Limit(limit).
		Find(&urls).Error
	return urls, err
}

func (r *urlRepository) GetRecentActivity(userID uuid.UUID, limit int) ([]domain.URL, error) {
	var urls []domain.URL
	err := r.db.Where("user_id = ? AND last_clicked_at IS NOT NULL", userID).
		Order("last_clicked_at DESC").
		Limit(limit).
		Find(&urls).Error
	return urls, err
}
