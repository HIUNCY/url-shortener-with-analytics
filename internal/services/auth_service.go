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

// NewAuthService membuat instance baru dari authService.
func NewAuthService(userRepo domain.UserRepository, cfg configs.Config) AuthService {
	return &authService{userRepo: userRepo, cfg: cfg}
}

func (s *authService) Register(req request.RegisterRequest) (*domain.User, error) {
	// 1. Cek apakah email sudah ada
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		// Jika err nil, berarti user ditemukan, email sudah ada.
		return nil, errors.New("AUTH_EMAIL_ALREADY_EXISTS")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Error lain selain record not found
		return nil, err
	}

	// 2. Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 3. Generate API Key
	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		return nil, err
	}

	// 4. Buat objek domain User baru
	newUser := &domain.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		APIKey:       apiKey,
		FirstName:    &req.FirstName,
		LastName:     &req.LastName,
		PlanType:     "free", // Default plan
	}

	// 5. Simpan user ke database
	if err := s.userRepo.Store(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *authService) Login(req request.LoginRequest) (*LoginResult, error) {
	// 1. Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("AUTH_INVALID_CREDENTIALS")
		}
		return nil, err
	}

	// 2. Verifikasi password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("AUTH_INVALID_CREDENTIALS")
	}

	// 3. Generate Access Token
	accessExpiresIn, _ := time.ParseDuration(s.cfg.JWT.ExpiresIn)
	accessToken, err := utils.GenerateToken(user.ID, s.cfg.JWT.SecretKey, accessExpiresIn)
	if err != nil {
		return nil, err
	}

	// 4. Generate Refresh Token
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
	// 1. Validasi refresh token
	claims, err := utils.ValidateToken(refreshToken, s.cfg.JWT.RefreshSecretKey)
	if err != nil {
		return "", errors.New("AUTH_INVALID_REFRESH_TOKEN")
	}

	// 2. Cek apakah user masih ada (opsional tapi bagus)
	_, err = s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return "", errors.New("AUTH_USER_NOT_FOUND")
	}

	// 3. Generate access token baru
	accessExpiresIn, _ := time.ParseDuration(s.cfg.JWT.ExpiresIn)
	newAccessToken, err := utils.GenerateToken(claims.UserID, s.cfg.JWT.SecretKey, accessExpiresIn)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

func (s *authService) Logout(accessToken string) error {
	// Di implementasi sederhana, kita tidak melakukan apa-apa di server.
	// Logika untuk mem-blacklist token akan ditambahkan di sini nanti saat mengintegrasikan Redis.
	return nil
}
