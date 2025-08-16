package services

import (
	"errors"
	"log"
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RedirectService interface {
	ProcessRedirect(c *gin.Context, shortCode string) (string, error)
}

type redirectService struct {
	urlRepo   domain.URLRepository
	clickRepo domain.ClickRepository
}

func NewRedirectService(urlRepo domain.URLRepository, clickRepo domain.ClickRepository) RedirectService {
	return &redirectService{urlRepo: urlRepo, clickRepo: clickRepo}
}

func (s *redirectService) ProcessRedirect(c *gin.Context, shortCode string) (string, error) {
	// 1. Cari URL berdasarkan shortCode
	url, err := s.urlRepo.FindByShortCode(shortCode)
	if err != nil {
		return "", errors.New("URL_NOT_FOUND")
	}

	// 2. Periksa status URL
	if !url.IsActive {
		return "", errors.New("URL_NOT_FOUND")
	}
	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		return "", errors.New("URL_NOT_FOUND")
	}
	if url.PasswordHash != nil {
		// Untuk sekarang, kita kembalikan error. Nanti ini bisa mengarah ke halaman password.
		return "", errors.New("URL_PASSWORD_PROTECTED")
	}

	// 3. Jalankan pelacakan klik secara asynchronous (agar tidak memperlambat redirect)
	go s.trackClick(c, url.ID)

	// 4. Kembalikan URL asli untuk di-redirect
	return url.OriginalURL, nil
}

// trackClick berjalan di background
func (s *redirectService) trackClick(c *gin.Context, urlID uuid.UUID) {
	// Update click count di tabel URLs
	if err := s.urlRepo.IncrementClickCount(urlID); err != nil {
		log.Printf("Error incrementing click count for URL %s: %v", urlID, err)
	}

	// Simpan detail klik di tabel clicks
	newClick := &domain.Click{
		URLID:      urlID,
		IPAddress:  c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Referer:    c.Request.Referer(),
		DeviceType: "unknown",
	}
	if err := s.clickRepo.Store(newClick); err != nil {
		log.Printf("Error storing click details for URL %s: %v", urlID, err)
	}
}
