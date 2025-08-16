package services

import (
	"errors"
	"log"
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UnlockResult struct {
	RedirectURL string
	AccessToken string
}

type RedirectService interface {
	ProcessRedirect(c *gin.Context, shortCode string) (string, error)
	UnlockURL(shortCode, password string) (*UnlockResult, error)
}

type redirectService struct {
	urlRepo   domain.URLRepository
	clickRepo domain.ClickRepository
	cfg       configs.Config
}

func NewRedirectService(urlRepo domain.URLRepository, clickRepo domain.ClickRepository, cfg configs.Config) RedirectService {
	return &redirectService{urlRepo: urlRepo, clickRepo: clickRepo, cfg: cfg}
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

func (s *redirectService) UnlockURL(shortCode, password string) (*UnlockResult, error) {
	// 1. Cari URL
	url, err := s.urlRepo.FindByShortCode(shortCode)
	if err != nil {
		return nil, errors.New("URL_NOT_FOUND")
	}

	// 2. Cek apakah URL punya password
	if url.PasswordHash == nil {
		return nil, errors.New("URL_NOT_PROTECTED")
	}

	// 3. Verifikasi password
	if !utils.CheckPasswordHash(password, *url.PasswordHash) {
		return nil, errors.New("URL_INVALID_PASSWORD")
	}

	// 4. Buat token redirect sementara (berlaku 1 menit)
	// Kita bisa gunakan kembali user ID dari URL, atau ID URL itu sendiri sebagai subjek.
	// Di sini kita gunakan ID URL untuk membuatnya spesifik.
	tempToken, err := utils.GenerateToken(url.ID, s.cfg.JWT.SecretKey, 1*time.Minute)
	if err != nil {
		return nil, err
	}

	return &UnlockResult{
		RedirectURL: url.OriginalURL,
		AccessToken: tempToken,
	}, nil
}
