package routes

import (
	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/handlers"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func SetupURLRoutes(router *gin.RouterGroup, urlHandler *handlers.URLHandler, cfg configs.Config, userRepo domain.UserRepository) {
	urlGroup := router.Group("/urls")
	urlGroup.Use(middleware.AuthMiddleware(cfg.JWT, userRepo))
	{
		urlGroup.POST("", urlHandler.CreateShortURL)
		urlGroup.GET("", urlHandler.GetUserURLs)
		urlGroup.GET("/:urlID", urlHandler.GetURLDetails)
		urlGroup.PUT("/:urlID", urlHandler.UpdateURL)
		urlGroup.DELETE("/:urlID", urlHandler.DeleteURL)
	}
}
