package handlers

import (
	"net/http"
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/response"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnalyticsHandler struct {
	analyticsService services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

// GetURLAnalytics godoc
// @Summary Get URL analytics
// @Description Retrieves detailed analytics for a specific URL.
// @Tags Analytics
// @Security BearerAuth
// @Produce  json
// @Param    url_id path string true "URL ID" format(uuid)
// @Param period query string false "Time period for analytics" Enums(24h, 7d, 30d, all) default(7d)
// @Success 200 {object} response.URLAnalyticsSuccessResponse
// @Failure 403 {object} response.APIErrorResponse "Forbidden"
// @Failure 404 {object} response.APIErrorResponse "URL not found"
// @Router /urls/{url_id}/analytics [get]
func (h *AnalyticsHandler) GetURLAnalytics(c *gin.Context) {
	urlID, _ := uuid.Parse(c.Param("urlID"))
	userID := c.MustGet("userID").(uuid.UUID)
	period := c.DefaultQuery("period", "7d")

	analyticsData, err := h.analyticsService.GetURLAnalytics(urlID, userID, period)
	if err != nil {
		// Handle not found, forbidden, etc.
		response.SendError(c, http.StatusNotFound, "NOT_FOUND", "Could not retrieve analytics for URL", nil)
		return
	}

	c.JSON(http.StatusOK, response.URLAnalyticsSuccessResponse{
		Success:   true,
		Data:      *analyticsData,
		Timestamp: time.Now().UTC(),
	})
}
