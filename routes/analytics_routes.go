package routes

import (
	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/handlers"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAnalyticsRoutes(router *gin.RouterGroup, analyticsHandler *handlers.AnalyticsHandler, cfg configs.Config, userRepo domain.UserRepository) {
	// Rute ini secara teknis berada di bawah URL, tetapi kita pisahkan untuk kerapian
	analyticsGroup := router.Group("/urls/:urlID/analytics")
	analyticsGroup.Use(middleware.AuthMiddleware(cfg.JWT, userRepo))
	{
		analyticsGroup.GET("", analyticsHandler.GetURLAnalytics)
	}
}
