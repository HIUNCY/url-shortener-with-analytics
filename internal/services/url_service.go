package services

import (
	"errors"
	"fmt"

	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/request"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateURLResult struct {
	URL      *domain.URL
	QRCode   string
	ShortURL string
}

type URLService interface {
	CreateShortURL(userID uuid.UUID, req request.CreateURLRequest) (*CreateURLResult, error)
}

type urlService struct {
	urlRepo domain.URLRepository
}

func NewURLService(urlRepo domain.URLRepository) URLService {
	return &urlService{urlRepo: urlRepo}
}

func (s *urlService) CreateShortURL(userID uuid.UUID, req request.CreateURLRequest) (*CreateURLResult, error) {
	shortCode := ""
	if req.CustomAlias != nil && *req.CustomAlias != "" {
		// Cek ketersediaan custom alias
		_, err := s.urlRepo.FindByCustomAlias(*req.CustomAlias)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("URL_CUSTOM_ALIAS_EXISTS")
		}
		shortCode = *req.CustomAlias
	} else {
		// Generate short code acak dan pastikan unik
		for {
			newCode, err := utils.GenerateShortCode()
			if err != nil {
				return nil, err
			}
			_, err = s.urlRepo.FindByShortCode(newCode)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				shortCode = newCode
				break
			}
		}
	}

	var hashedPassword *string
	if req.Password != nil && *req.Password != "" {
		hash, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		hashedPassword = &hash
	}

	newURL := &domain.URL{
		UserID:       &userID,
		OriginalURL:  req.OriginalURL,
		ShortCode:    shortCode,
		CustomAlias:  req.CustomAlias,
		Title:        req.Title,
		Description:  req.Description,
		ExpiresAt:    req.ExpiresAt,
		PasswordHash: hashedPassword,
	}

	if err := s.urlRepo.Store(newURL); err != nil {
		return nil, err
	}

	// TODO: Ganti "http://localhost:8080" dengan domain dari config
	shortURLString := fmt.Sprintf("http://localhost:8080/%s", newURL.ShortCode)
	qrCode, err := utils.GenerateQRCodeBase64(shortURLString, 256)
	if err != nil {
		// Log error tapi jangan gagalkan proses utama
		fmt.Printf("Gagal generate QR Code untuk URL %s: %v\n", newURL.ID, err)
	}

	return &CreateURLResult{
		URL:      newURL,
		QRCode:   qrCode,
		ShortURL: shortURLString,
	}, nil
}
