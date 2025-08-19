package routes

import (
	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/handlers"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func SetupQRCodeRoutes(router *gin.RouterGroup, qrCodeHandler *handlers.QRCodeHandler, cfg configs.Config, userRepo domain.UserRepository) {
	qrGroup := router.Group("/urls/:urlID/qr")
	qrGroup.Use(middleware.AuthMiddleware(cfg.JWT, userRepo))
	{
		qrGroup.GET("", qrCodeHandler.GetQRCode)
		qrGroup.GET("/download", qrCodeHandler.DownloadQRCode)
	}
}
