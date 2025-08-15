package middleware

import (
	"net/http"
	"strings"

	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/response"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware(cfg configs.JWTConfig, userRepo domain.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userID uuid.UUID

		// Coba autentikasi via Bearer Token
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := utils.ValidateToken(tokenString, cfg.SecretKey)
			if err != nil {
				response.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token", nil)
				return
			}
			userID = claims.UserID
		} else {
			// Jika tidak ada Bearer Token, coba via API Key
			apiKey := c.GetHeader("X-API-Key")
			if apiKey == "" {
				response.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header or X-API-Key header is required", nil)
				return
			}
			user, err := userRepo.FindByAPIKey(apiKey)
			if err != nil {
				response.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid API Key", nil)
				return
			}
			userID = user.ID
		}

		if userID == uuid.Nil {
			response.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Could not authenticate user", nil)
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
