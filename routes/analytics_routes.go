package routes

import (
	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/handlers"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAnalyticsRoutes(router *gin.RouterGroup, analyticsHandler *handlers.AnalyticsHandler, cfg configs.Config, userRepo domain.UserRepository) {
	analyticsGroup := router.Group("/analytics")
	analyticsGroup.Use(middleware.AuthMiddleware(cfg.JWT, userRepo))
	{
		analyticsGroup.GET("/dashboard", analyticsHandler.GetUserDashboard)
	}

	urlAnalyticsGroup := router.Group("/urls/:urlID/analytics")
	urlAnalyticsGroup.Use(middleware.AuthMiddleware(cfg.JWT, userRepo))
	{
		urlAnalyticsGroup.GET("", analyticsHandler.GetURLAnalytics)
	}
}
