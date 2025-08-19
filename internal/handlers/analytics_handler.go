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
		if err.Error() == "URL_FORBIDDEN" {
			response.SendError(c, http.StatusForbidden, "FORBIDDEN", "You do not have permission to view this URL", nil)
			return
		}
		response.SendError(c, http.StatusNotFound, "NOT_FOUND", "Could not retrieve analytics for URL", nil)
		return
	}

	c.JSON(http.StatusOK, response.URLAnalyticsSuccessResponse{
		Success:   true,
		Data:      *analyticsData,
		Timestamp: time.Now().UTC(),
	})
}

// GetUserDashboard godoc
// @Summary Get user dashboard analytics
// @Description Retrieves summary analytics for the authenticated user's dashboard.
// @Tags Analytics
// @Security BearerAuth
// @Produce  json
// @Success 200 {object} response.UserDashboardSuccessResponse
// @Failure 401 {object} response.APIErrorResponse "Unauthorized"
// @Router /analytics/dashboard [get]
func (h *AnalyticsHandler) GetUserDashboard(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	dashboardData, err := h.analyticsService.GetUserDashboard(userID)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to retrieve dashboard data", nil)
		return
	}

	c.JSON(http.StatusOK, response.UserDashboardSuccessResponse{
		Success:   true,
		Data:      *dashboardData,
		Timestamp: time.Now().UTC(),
	})
}
