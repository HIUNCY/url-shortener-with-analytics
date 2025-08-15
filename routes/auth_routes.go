package routes

import (
	"github.com/HIUNCY/url-shortener-with-analytics/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes mengatur rute untuk autentikasi.
func SetupAuthRoutes(router *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}
}
