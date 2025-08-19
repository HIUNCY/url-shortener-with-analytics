package response

import (
	"time"

	"github.com/google/uuid"
)

type TimeSeriesStat struct {
	Date   string `json:"date"`
	Clicks int64  `json:"clicks"`
}

type GroupedStat struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

type AnalyticsOverview struct {
	TotalClicks int64  `json:"total_clicks"`
	TopReferrer string `json:"top_referrer"`
	TopCountry  string `json:"top_country"`
}

type URLAnalyticsResponse struct {
	Overview         AnalyticsOverview `json:"overview"`
	ClicksOverTime   []TimeSeriesStat  `json:"clicks_over_time"`
	Referrers        []GroupedStat     `json:"referrers"`
	Countries        []GroupedStat     `json:"countries"`
	Devices          []GroupedStat     `json:"devices"`
	Browsers         []GroupedStat     `json:"browsers"`
	OperatingSystems []GroupedStat     `json:"operating_systems"`
}

type URLAnalyticsSuccessResponse struct {
	Success   bool                 `json:"success" example:"true"`
	Data      URLAnalyticsResponse `json:"data"`
	Timestamp time.Time            `json:"timestamp"`
}

type DashboardSummary struct {
	TotalURLs   int64 `json:"total_urls"`
	TotalClicks int64 `json:"total_clicks"`
	ActiveURLs  int64 `json:"active_urls"`
}

type DashboardActivityItem struct {
	URLID         uuid.UUID  `json:"url_id"`
	ShortCode     string     `json:"short_code"`
	Title         *string    `json:"title,omitempty"`
	LastClickedAt *time.Time `json:"last_clicked_at,omitempty"`
}

type DashboardTopURL struct {
	URLID      uuid.UUID `json:"url_id"`
	ShortCode  string    `json:"short_code"`
	Title      *string   `json:"title,omitempty"`
	ClickCount int       `json:"click_count"`
}

type UserDashboardResponse struct {
	Summary           DashboardSummary        `json:"summary"`
	RecentActivity    []DashboardActivityItem `json:"recent_activity"`
	TopPerformingURLs []DashboardTopURL       `json:"top_performing_urls"`
}

type UserDashboardSuccessResponse struct {
	Success   bool                  `json:"success" example:"true"`
	Data      UserDashboardResponse `json:"data"`
	Timestamp time.Time             `json:"timestamp"`
}
