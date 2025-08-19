package services

import (
	"errors"
	"fmt"

	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/request"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/response"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateURLResult struct {
	URL      *domain.URL
	QRCode   string
	ShortURL string
}

type URLListResult struct {
	URLs       []domain.URL
	Pagination response.PaginationResponse
}

type URLService interface {
	CreateShortURL(userID uuid.UUID, req request.CreateURLRequest) (*CreateURLResult, error)
	GetURLDetails(urlID, userID uuid.UUID) (*domain.URL, error)
	GetUserURLs(userID uuid.UUID, options *domain.FindAllOptions) (*URLListResult, error)
	UpdateURL(urlID, userID uuid.UUID, req request.UpdateURLRequest) (*domain.URL, error)
	DeleteURL(urlID, userID uuid.UUID) error
}

type urlService struct {
	urlRepo domain.URLRepository
	cfg     configs.Config
}

func NewURLService(urlRepo domain.URLRepository, cfg configs.Config) URLService {
	return &urlService{urlRepo: urlRepo, cfg: cfg}
}

func (s *urlService) CreateShortURL(userID uuid.UUID, req request.CreateURLRequest) (*CreateURLResult, error) {
	shortCode := ""
	if req.CustomAlias != nil && *req.CustomAlias != "" {
		_, err := s.urlRepo.FindByCustomAlias(*req.CustomAlias)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("URL_CUSTOM_ALIAS_EXISTS")
		}
		shortCode = *req.CustomAlias
	} else {
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

	shortURLString := fmt.Sprintf("%s/%s", s.cfg.Server.BaseURL, newURL.ShortCode)
	qrCode, err := utils.GenerateQRCodeBase64(shortURLString, 256)
	if err != nil {
		fmt.Printf("Gagal generate QR Code untuk URL %s: %v\n", newURL.ID, err)
	}

	return &CreateURLResult{
		URL:      newURL,
		QRCode:   qrCode,
		ShortURL: shortURLString,
	}, nil
}

func (s *urlService) GetUserURLs(userID uuid.UUID, options *domain.FindAllOptions) (*URLListResult, error) {
	urls, total, err := s.urlRepo.FindAllByUserID(userID, options)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if options.Limit > 0 {
		totalPages = int((total + int64(options.Limit) - 1) / int64(options.Limit))
	}

	pagination := response.PaginationResponse{
		Page:       (options.Offset / options.Limit) + 1,
		Limit:      options.Limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return &URLListResult{
		URLs:       urls,
		Pagination: pagination,
	}, nil
}

func (s *urlService) GetURLDetails(urlID, userID uuid.UUID) (*domain.URL, error) {
	url, err := s.urlRepo.FindByID(urlID)
	if err != nil {
		return nil, err
	}

	if url.UserID == nil || *url.UserID != userID {
		return nil, errors.New("URL_FORBIDDEN")
	}

	return url, nil
}

func (s *urlService) UpdateURL(urlID, userID uuid.UUID, req request.UpdateURLRequest) (*domain.URL, error) {
	url, err := s.urlRepo.FindByID(urlID)
	if err != nil {
		return nil, err
	}
	if url.UserID == nil || *url.UserID != userID {
		return nil, errors.New("URL_FORBIDDEN")
	}

	if req.Title != nil {
		url.Title = req.Title
	}
	if req.Description != nil {
		url.Description = req.Description
	}
	if req.ExpiresAt != nil {
		url.ExpiresAt = req.ExpiresAt
	}
	if req.IsActive != nil {
		url.IsActive = *req.IsActive
	}

	if err := s.urlRepo.Update(url); err != nil {
		return nil, err
	}
	return url, nil
}

func (s *urlService) DeleteURL(urlID, userID uuid.UUID) error {
	url, err := s.urlRepo.FindByID(urlID)
	if err != nil {
		return err
	}
	if url.UserID == nil || *url.UserID != userID {
		return errors.New("URL_FORBIDDEN")
	}

	return s.urlRepo.Delete(url)
}
