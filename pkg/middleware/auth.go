package middleware

import (
	"net/http"
	"strings"

	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/response"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware adalah middleware untuk memvalidasi JWT Access Token.
func AuthMiddleware(cfg configs.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header is required", nil)
			return
		}

		// Format header harus "Bearer {token}"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header format must be Bearer {token}", nil)
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString, cfg.SecretKey)
		if err != nil {
			response.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token", nil)
			return
		}

		// Simpan user ID di context untuk digunakan oleh handler selanjutnya
		c.Set("userID", claims.UserID)

		c.Next()
	}
}
