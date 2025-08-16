package main

import (
	"fmt"
	"log"

	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	_ "github.com/HIUNCY/url-shortener-with-analytics/docs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/handlers"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/repository/postgres"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/services"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/database"
	"github.com/HIUNCY/url-shortener-with-analytics/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title URL Shortener API
// @version 1.0
// @description This is a URL shortener service with analytics.
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Memuat konfigurasi
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("Tidak dapat memuat konfigurasi: %v", err)
	}

	// 2. Membuat koneksi database
	db, err := database.NewPostgresConnection(&config.Database)
	if err != nil {
		return
	}

	// 3. Inisialisasi semua komponen (Dependency Injection)
	userRepository := postgres.NewUserRepository(db)
	urlRepository := postgres.NewURLRepository(db)
	clickRepository := postgres.NewClickRepository(db)
	authService := services.NewAuthService(userRepository, config)
	userService := services.NewUserService(userRepository)
	urlService := services.NewURLService(urlRepository, config)
	redirectService := services.NewRedirectService(urlRepository, clickRepository, config)
	authHandler := handlers.NewAuthHandler(authService, config)
	profileHandler := handlers.NewProfileHandler(userService)
	urlHandler := handlers.NewURLHandler(urlService, config)
	redirectHandler := handlers.NewRedirectHandler(redirectService, config)

	// 4. Setup Gin Router
	router := gin.Default()

	router.GET("/:shortCode", redirectHandler.Redirect)
	router.POST("/:shortCode/unlock", redirectHandler.UnlockURL)
	router.GET("/:shortCode/info", redirectHandler.GetURLInfo)

	// Route untuk Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Grup rute untuk API v1
	apiV1 := router.Group("/api/v1")
	routes.SetupAuthRoutes(apiV1, authHandler, config, userRepository)
	routes.SetupProfileRoutes(apiV1, profileHandler, config, userRepository)
	routes.SetupURLRoutes(apiV1, urlHandler, config, userRepository)

	// 5. Jalankan server
	serverAddress := fmt.Sprintf(":%s", config.Server.Port)
	log.Printf("Server berjalan di %s", serverAddress)
	if err := router.Run(serverAddress); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
