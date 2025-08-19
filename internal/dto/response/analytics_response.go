package response

import "time"

// TimeSeriesStat adalah statistik untuk data berbasis waktu.
type TimeSeriesStat struct {
	Date   string `json:"date"`
	Clicks int64  `json:"clicks"`
}

// GroupedStat adalah statistik umum untuk data yang dikelompokkan.
type GroupedStat struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

// AnalyticsOverview berisi ringkasan data analitik.
type AnalyticsOverview struct {
	TotalClicks int64  `json:"total_clicks"`
	TopReferrer string `json:"top_referrer"`
	TopCountry  string `json:"top_country"`
}

// URLAnalyticsResponse adalah DTO untuk payload data pada respons analitik URL.
type URLAnalyticsResponse struct {
	Overview         AnalyticsOverview `json:"overview"`
	ClicksOverTime   []TimeSeriesStat  `json:"clicks_over_time"`
	Referrers        []GroupedStat     `json:"referrers"`
	Countries        []GroupedStat     `json:"countries"`
	Devices          []GroupedStat     `json:"devices"`
	Browsers         []GroupedStat     `json:"browsers"`
	OperatingSystems []GroupedStat     `json:"operating_systems"`
}

// URLAnalyticsSuccessResponse adalah wrapper untuk Swagger.
type URLAnalyticsSuccessResponse struct {
	Success   bool                 `json:"success" example:"true"`
	Data      URLAnalyticsResponse `json:"data"`
	Timestamp time.Time            `json:"timestamp"`
}
