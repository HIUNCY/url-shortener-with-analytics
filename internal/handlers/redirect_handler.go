package handlers

import (
	"net/http"

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
