package services

import (
	"errors"
	"fmt"

	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/utils"
	"github.com/google/uuid"
)

type QRCodeService interface {
	GetQRCodeInfo(urlID, userID uuid.UUID, size int, format string) (*domain.URL, string, error)
	GetQRCodeForDownload(urlID, userID uuid.UUID, size int) (*domain.URL, []byte, error)
}

type qrCodeService struct {
	urlRepo domain.URLRepository
	cfg     configs.Config
}

func NewQRCodeService(urlRepo domain.URLRepository, cfg configs.Config) QRCodeService {
	return &qrCodeService{urlRepo: urlRepo, cfg: cfg}
}

// getAndVerifyURL adalah fungsi internal untuk menghindari duplikasi kode.
func (s *qrCodeService) getAndVerifyURL(urlID, userID uuid.UUID) (*domain.URL, error) {
	url, err := s.urlRepo.FindByID(urlID)
	if err != nil {
		return nil, errors.New("URL_NOT_FOUND")
	}
	if url.UserID == nil || *url.UserID != userID {
		return nil, errors.New("URL_FORBIDDEN")
	}
	return url, nil
}

func (s *qrCodeService) GetQRCodeInfo(urlID, userID uuid.UUID, size int, format string) (*domain.URL, string, error) {
	url, err := s.getAndVerifyURL(urlID, userID)
	if err != nil {
		return nil, "", err
	}

	shortURL := fmt.Sprintf("%s/%s", s.cfg.Server.BaseURL, url.ShortCode)
	qrCodeBase64, err := utils.GenerateQRCodeBase64(shortURL, size)
	if err != nil {
		return nil, "", err
	}

	return url, qrCodeBase64, nil
}

func (s *qrCodeService) GetQRCodeForDownload(urlID, userID uuid.UUID, size int) (*domain.URL, []byte, error) {
	url, err := s.getAndVerifyURL(urlID, userID)
	if err != nil {
		return nil, nil, err
	}

	shortURL := fmt.Sprintf("%s/%s", s.cfg.Server.BaseURL, url.ShortCode)
	qrCodeBytes, err := utils.GenerateQRCodeBytes(shortURL, size)
	if err != nil {
		return nil, nil, err
	}

	return url, qrCodeBytes, nil
}
