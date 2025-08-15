package routes

import (
	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/handlers"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes mengatur rute untuk autentikasi.
func SetupAuthRoutes(router *gin.RouterGroup, authHandler *handlers.AuthHandler, cfg configs.Config) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)

		protected := authGroup.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT))
		{
			protected.POST("/logout", authHandler.Logout)
		}
	}
}
