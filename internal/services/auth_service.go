package services

import (
	"errors"

	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/request"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/utils"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req request.RegisterRequest) (*domain.User, error)
}

type authService struct {
	userRepo domain.UserRepository
}

// NewAuthService membuat instance baru dari authService.
func NewAuthService(userRepo domain.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
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
