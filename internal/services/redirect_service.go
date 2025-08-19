package services

import (
	"errors"
	"log"
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/geoip"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UnlockResult struct {
	RedirectURL string
	AccessToken string
}

type InfoResult struct {
	URL    *domain.URL
	Domain string
	IsSafe bool
}

type RedirectService interface {
	ProcessRedirect(c *gin.Context, shortCode string) (string, error)
	UnlockURL(shortCode, password string) (*UnlockResult, error)
	GetURLInfo(shortCode string) (*InfoResult, error)
}

type redirectService struct {
	urlRepo   domain.URLRepository
	clickRepo domain.ClickRepository
	geoipSvc  geoip.GeoIPService
	cfg       configs.Config
}

func NewRedirectService(urlRepo domain.URLRepository, clickRepo domain.ClickRepository, geoipSvc geoip.GeoIPService, cfg configs.Config) RedirectService {
	return &redirectService{urlRepo: urlRepo, clickRepo: clickRepo, geoipSvc: geoipSvc, cfg: cfg}
}

func (s *redirectService) ProcessRedirect(c *gin.Context, shortCode string) (string, error) {
	url, err := s.urlRepo.FindByShortCode(shortCode)
	if err != nil {
		return "", errors.New("URL_NOT_FOUND")
	}

	if !url.IsActive {
		return "", errors.New("URL_NOT_FOUND")
	}
	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		return "", errors.New("URL_NOT_FOUND")
	}
	if url.PasswordHash != nil {
		return "", errors.New("URL_PASSWORD_PROTECTED")
	}

	go s.trackClick(c, url.ID)

	return url.OriginalURL, nil
}

func (s *redirectService) trackClick(c *gin.Context, urlID uuid.UUID) {
	if err := s.urlRepo.IncrementClickCount(urlID); err != nil {
		log.Printf("Error incrementing click count for URL %s: %v", urlID, err)
	}

	uaString := c.Request.UserAgent()
	parsedUA := utils.ParseUserAgent(uaString)
	clientIP := c.ClientIP()

	location, err := s.geoipSvc.Lookup(clientIP)
	if err != nil {
		log.Printf("Could not perform GeoIP lookup for IP %s: %v", clientIP, err)
	}

	newClick := &domain.Click{
		URLID:      urlID,
		IPAddress:  clientIP,
		UserAgent:  uaString,
		Referer:    c.Request.Referer(),
		DeviceType: parsedUA.DeviceType,
		Browser:    parsedUA.BrowserName,
		OS:         parsedUA.OSName,
		Country:    location.Country,
		Region:     location.Region,
		City:       location.City,
	}
	if err := s.clickRepo.Store(newClick); err != nil {
		log.Printf("Error storing click details for URL %s: %v", urlID, err)
	}
}

func (s *redirectService) UnlockURL(shortCode, password string) (*UnlockResult, error) {
	url, err := s.urlRepo.FindByShortCode(shortCode)
	if err != nil {
		return nil, errors.New("URL_NOT_FOUND")
	}

	if url.PasswordHash == nil {
		return nil, errors.New("URL_NOT_PROTECTED")
	}

	if !utils.CheckPasswordHash(password, *url.PasswordHash) {
		return nil, errors.New("URL_INVALID_PASSWORD")
	}

	tempToken, err := utils.GenerateToken(url.ID, s.cfg.JWT.SecretKey, 1*time.Minute)
	if err != nil {
		return nil, err
	}

	return &UnlockResult{
		RedirectURL: url.OriginalURL,
		AccessToken: tempToken,
	}, nil
}

func (s *redirectService) GetURLInfo(shortCode string) (*InfoResult, error) {
	url, err := s.urlRepo.FindByShortCode(shortCode)
	if err != nil {
		return nil, errors.New("URL_NOT_FOUND")
	}

	if !url.IsActive || (url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now())) {
		return nil, errors.New("URL_NOT_FOUND")
	}

	domainName, err := utils.GetDomainFromURL(url.OriginalURL)
	if err != nil {
		domainName = ""
	}

	isSafe := true

	return &InfoResult{
		URL:    url,
		Domain: domainName,
		IsSafe: isSafe,
	}, nil
}
