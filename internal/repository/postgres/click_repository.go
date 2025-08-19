package postgres

import (
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/google/uuid"
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

func (r *clickRepository) getAggregatedStats(urlID uuid.UUID, since time.Time, limit int, column string) ([]domain.GroupedResult, error) {
	var results []domain.GroupedResult
	err := r.db.Model(&domain.Click{}).
		Select(column+" as value, COUNT(*) as count").
		Where("url_id = ? AND clicked_at >= ?", urlID, since).
		Where(column + " IS NOT NULL AND " + column + " != ''").
		Group(column).
		Order("count DESC").
		Limit(limit).
		Find(&results).Error
	return results, err
}

func (r *clickRepository) GetTotalClicks(urlID uuid.UUID, since time.Time) (int64, error) {
	var total int64
	err := r.db.Model(&domain.Click{}).Where("url_id = ? AND clicked_at >= ?", urlID, since).Count(&total).Error
	return total, err
}

func (r *clickRepository) GetTopReferrer(urlID uuid.UUID, since time.Time) (string, error) {
	var result domain.GroupedResult
	err := r.db.Model(&domain.Click{}).Select("referer as value, COUNT(*) as count").
		Where("url_id = ? AND clicked_at >= ? AND referer IS NOT NULL AND referer != ''", urlID, since).
		Group("referer").Order("count DESC").First(&result).Error
	return result.Value, err
}
func (r *clickRepository) GetTopCountry(urlID uuid.UUID, since time.Time) (string, error) {
	var result domain.GroupedResult
	err := r.db.Model(&domain.Click{}).Select("country as value, COUNT(*) as count").
		Where("url_id = ? AND clicked_at >= ? AND country IS NOT NULL AND country != ''", urlID, since).
		Group("country").Order("count DESC").First(&result).Error
	return result.Value, err
}
func (r *clickRepository) GetClicksOverTime(urlID uuid.UUID, since time.Time) ([]domain.TimeSeriesResult, error) {
	var results []domain.TimeSeriesResult
	err := r.db.Model(&domain.Click{}).
		Select("DATE(clicked_at) as date, COUNT(*) as count").
		Where("url_id = ? AND clicked_at >= ?", urlID, since).
		Group("DATE(clicked_at)").
		Order("date ASC").
		Find(&results).Error
	return results, err
}

func (r *clickRepository) GetTopCountries(urlID uuid.UUID, since time.Time, limit int) ([]domain.GroupedResult, error) {
	return r.getAggregatedStats(urlID, since, limit, "country")
}
func (r *clickRepository) GetTopReferrers(urlID uuid.UUID, since time.Time, limit int) ([]domain.GroupedResult, error) {
	return r.getAggregatedStats(urlID, since, limit, "referer")
}
func (r *clickRepository) GetDeviceStats(urlID uuid.UUID, since time.Time) ([]domain.GroupedResult, error) {
	return r.getAggregatedStats(urlID, since, 10, "device_type")
}
func (r *clickRepository) GetBrowserStats(urlID uuid.UUID, since time.Time) ([]domain.GroupedResult, error) {
	return r.getAggregatedStats(urlID, since, 10, "browser")
}
func (r *clickRepository) GetOSStats(urlID uuid.UUID, since time.Time) ([]domain.GroupedResult, error) {
	return r.getAggregatedStats(urlID, since, 10, "os")
}
