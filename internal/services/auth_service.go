package services

import (
	"errors"
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/request"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/utils"
	"gorm.io/gorm"
)

type LoginResult struct {
	User         *domain.User
	AccessToken  string
	RefreshToken string
}

type AuthService interface {
	Register(req request.RegisterRequest) (*domain.User, error)
	Login(req request.LoginRequest) (*LoginResult, error)
	RefreshToken(refreshToken string) (string, error)
	Logout(accessToken string) error
}

type authService struct {
	userRepo domain.UserRepository
	cfg      configs.Config
}

func NewAuthService(userRepo domain.UserRepository, cfg configs.Config) AuthService {
	return &authService{userRepo: userRepo, cfg: cfg}
}

func (s *authService) Register(req request.RegisterRequest) (*domain.User, error) {
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("AUTH_EMAIL_ALREADY_EXISTS")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		APIKey:       apiKey,
		FirstName:    &req.FirstName,
		LastName:     &req.LastName,
		PlanType:     "free",
	}

	if err := s.userRepo.Store(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *authService) Login(req request.LoginRequest) (*LoginResult, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("AUTH_INVALID_CREDENTIALS")
		}
		return nil, err
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("AUTH_INVALID_CREDENTIALS")
	}

	accessExpiresIn, _ := time.ParseDuration(s.cfg.JWT.ExpiresIn)
	accessToken, err := utils.GenerateToken(user.ID, s.cfg.JWT.SecretKey, accessExpiresIn)
	if err != nil {
		return nil, err
	}

	refreshExpiresIn, _ := time.ParseDuration(s.cfg.JWT.RefreshExpiresIn)
	refreshToken, err := utils.GenerateToken(user.ID, s.cfg.JWT.RefreshSecretKey, refreshExpiresIn)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshToken(refreshToken string) (string, error) {
	claims, err := utils.ValidateToken(refreshToken, s.cfg.JWT.RefreshSecretKey)
	if err != nil {
		return "", errors.New("AUTH_INVALID_REFRESH_TOKEN")
	}

	_, err = s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return "", errors.New("AUTH_USER_NOT_FOUND")
	}

	accessExpiresIn, _ := time.ParseDuration(s.cfg.JWT.ExpiresIn)
	newAccessToken, err := utils.GenerateToken(claims.UserID, s.cfg.JWT.SecretKey, accessExpiresIn)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

func (s *authService) Logout(accessToken string) error {
	return nil
}
