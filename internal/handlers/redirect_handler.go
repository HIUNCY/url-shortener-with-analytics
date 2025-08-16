package handlers

import (
	"net/http"
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/request"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/response"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/services"
	"github.com/gin-gonic/gin"
)

type RedirectHandler struct {
	redirectService services.RedirectService
}

func NewRedirectHandler(redirectService services.RedirectService) *RedirectHandler {
	return &RedirectHandler{redirectService: redirectService}
}

// Redirect menangani request ke short URL
func (h *RedirectHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("shortCode")

	originalURL, err := h.redirectService.ProcessRedirect(c, shortCode)
	if err != nil {
		if err.Error() == "URL_PASSWORD_PROTECTED" {
			// Nanti bisa redirect ke halaman input password
			response.SendError(c, http.StatusUnauthorized, "PASSWORD_PROTECTED", "This URL is password protected", nil)
			return
		}
		// Untuk semua error lain (not found, expired, inactive), kita tampilkan 404
		c.HTML(http.StatusNotFound, "404.html", nil) // Anggap kita punya template 404.html
		return
	}

	// Lakukan redirect
	c.Redirect(http.StatusFound, originalURL)
}

// UnlockURL godoc
// @Summary Unlock a password-protected URL
// @Description Verifies the password for a short URL and returns the original URL.
// @Tags Redirection
// @Accept   json
// @Produce  json
// @Param    shortCode path string true "Short Code"
// @Param    password body request.UnlockURLRequest true "Password"
// @Success 200 {object} response.UnlockURLSuccessResponse "URL unlocked successfully"
// @Failure 400 {object} response.APIErrorResponse "Validation error"
// @Failure 401 {object} response.APIErrorResponse "Invalid password"
// @Failure 404 {object} response.APIErrorResponse "URL not found"
// @Router /{shortCode}/unlock [post]
func (h *RedirectHandler) UnlockURL(c *gin.Context) {
	shortCode := c.Param("shortCode")
	var req request.UnlockURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	result, err := h.redirectService.UnlockURL(shortCode, req.Password)
	if err != nil {
		if err.Error() == "URL_INVALID_PASSWORD" {
			response.SendError(c, http.StatusUnauthorized, "INVALID_PASSWORD", "The provided password is incorrect", nil)
			return
		}
		// Handle error not found, dll.
		response.SendError(c, http.StatusNotFound, "NOT_FOUND", "URL not found or not password protected", nil)
		return
	}

	c.JSON(http.StatusOK, response.UnlockURLSuccessResponse{
		Success: true,
		Data: response.UnlockURLResponse{
			RedirectURL: result.RedirectURL,
			AccessToken: result.AccessToken,
		},
		Timestamp: time.Now().UTC(),
	})
}
