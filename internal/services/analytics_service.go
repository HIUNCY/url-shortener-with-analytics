package services

import (
	"errors"
	"sync"
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/response"
	"github.com/google/uuid"
)

type AnalyticsService interface {
	GetURLAnalytics(urlID, userID uuid.UUID, period string) (*response.URLAnalyticsResponse, error)
	GetUserDashboard(userID uuid.UUID) (*response.UserDashboardResponse, error)
}

type analyticsService struct {
	urlRepo   domain.URLRepository
	clickRepo domain.ClickRepository
}

func NewAnalyticsService(urlRepo domain.URLRepository, clickRepo domain.ClickRepository) AnalyticsService {
	return &analyticsService{urlRepo: urlRepo, clickRepo: clickRepo}
}

func (s *analyticsService) GetURLAnalytics(urlID, userID uuid.UUID, period string) (*response.URLAnalyticsResponse, error) {
	// 1. Verifikasi kepemilikan URL
	url, err := s.urlRepo.FindByID(urlID)
	if err != nil {
		return nil, errors.New("URL_NOT_FOUND")
	}
	if url.UserID == nil || *url.UserID != userID {
		return nil, errors.New("URL_FORBIDDEN")
	}

	// 2. Tentukan rentang waktu
	since := time.Now()
	switch period {
	case "24h":
		since = since.Add(-24 * time.Hour)
	case "7d":
		since = since.AddDate(0, 0, -7)
	case "30d":
		since = since.AddDate(0, 0, -30)
	default: // "all"
		since = time.Time{} // Waktu nol
	}

	// 3. Panggil semua query agregasi secara paralel
	var wg sync.WaitGroup
	var analyticsData response.URLAnalyticsResponse
	var errs = make(chan error, 10)

	// Overview
	wg.Add(3)
	go func() {
		defer wg.Done()
		analyticsData.Overview.TotalClicks, _ = s.clickRepo.GetTotalClicks(urlID, since)
	}()
	go func() {
		defer wg.Done()
		analyticsData.Overview.TopReferrer, _ = s.clickRepo.GetTopReferrer(urlID, since)
	}()
	go func() {
		defer wg.Done()
		analyticsData.Overview.TopCountry, _ = s.clickRepo.GetTopCountry(urlID, since)
	}()

	// Lists
	wg.Add(5)
	go func() {
		defer wg.Done()
		res, _ := s.clickRepo.GetClicksOverTime(urlID, since)
		analyticsData.ClicksOverTime = mapTimeSeries(res)
	}()
	go func() {
		defer wg.Done()
		res, _ := s.clickRepo.GetTopReferrers(urlID, since, 10)
		analyticsData.Referrers = mapGrouped(res)
	}()
	go func() {
		defer wg.Done()
		res, _ := s.clickRepo.GetTopCountries(urlID, since, 10)
		analyticsData.Countries = mapGrouped(res)
	}()
	go func() {
		defer wg.Done()
		res, _ := s.clickRepo.GetDeviceStats(urlID, since)
		analyticsData.Devices = mapGrouped(res)
	}()
	go func() {
		defer wg.Done()
		res, _ := s.clickRepo.GetBrowserStats(urlID, since)
		analyticsData.Browsers = mapGrouped(res)
	}()

	wg.Wait()
	close(errs)

	return &analyticsData, nil
}

// Helper untuk mapping
func mapTimeSeries(res []domain.TimeSeriesResult) []response.TimeSeriesStat {
	stats := make([]response.TimeSeriesStat, len(res))
	for i, r := range res {
		stats[i] = response.TimeSeriesStat{Date: r.Date.Format("2006-01-02"), Clicks: r.Count}
	}
	return stats
}
func mapGrouped(res []domain.GroupedResult) []response.GroupedStat {
	stats := make([]response.GroupedStat, len(res))
	for i, r := range res {
		stats[i] = response.GroupedStat{Value: r.Value, Count: r.Count}
	}
	return stats
}

func (s *analyticsService) GetUserDashboard(userID uuid.UUID) (*response.UserDashboardResponse, error) {
	var wg sync.WaitGroup
	dashboardData := &response.UserDashboardResponse{}

	// Panggil semua query secara paralel
	wg.Add(3)

	go func() {
		defer wg.Done()
		summary, _ := s.urlRepo.GetDashboardSummary(userID)
		if summary != nil {
			dashboardData.Summary = response.DashboardSummary(*summary)
		}
	}()

	go func() {
		defer wg.Done()
		topURLs, _ := s.urlRepo.GetTopPerformingURLs(userID, 5)
		dashboardData.TopPerformingURLs = make([]response.DashboardTopURL, len(topURLs))
		for i, u := range topURLs {
			dashboardData.TopPerformingURLs[i] = response.DashboardTopURL{
				URLID: u.ID, ShortCode: u.ShortCode, Title: u.Title, ClickCount: u.ClickCount,
			}
		}
	}()

	go func() {
		defer wg.Done()
		recentURLs, _ := s.urlRepo.GetRecentActivity(userID, 5)
		dashboardData.RecentActivity = make([]response.DashboardActivityItem, len(recentURLs))
		for i, u := range recentURLs {
			dashboardData.RecentActivity[i] = response.DashboardActivityItem{
				URLID: u.ID, ShortCode: u.ShortCode, Title: u.Title, LastClickedAt: u.LastClickedAt,
			}
		}
	}()

	wg.Wait()

	return dashboardData, nil
}
