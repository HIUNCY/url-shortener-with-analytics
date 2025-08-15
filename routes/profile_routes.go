package routes

import (
	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/handlers"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProfileRoutes(router *gin.RouterGroup, profileHandler *handlers.ProfileHandler, cfg configs.Config, userRepo domain.UserRepository) {
	profileGroup := router.Group("/profile")
	profileGroup.Use(middleware.AuthMiddleware(cfg.JWT, userRepo))
	{
		profileGroup.GET("", profileHandler.GetProfile)
		profileGroup.PUT("", profileHandler.UpdateProfile)
		profileGroup.PUT("/password", profileHandler.ChangePassword)
		profileGroup.POST("/api-key/regenerate", profileHandler.RegenerateAPIKey)
	}
}
