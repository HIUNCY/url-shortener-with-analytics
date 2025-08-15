package services

import (
	"errors"

	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/request"
	"github.com/HIUNCY/url-shortener-with-analytics/pkg/utils"
	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(userID uuid.UUID) (*domain.User, error)
	UpdateProfile(userID uuid.UUID, req request.UpdateProfileRequest) (*domain.User, error)
	ChangePassword(userID uuid.UUID, req request.ChangePasswordRequest) error
	RegenerateAPIKey(userID uuid.UUID) (string, error)
}

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetProfile(userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *userService) UpdateProfile(userID uuid.UUID, req request.UpdateProfileRequest) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	user.FirstName = &req.FirstName
	user.LastName = &req.LastName

	err = s.userRepo.Update(user)
	return user, err
}

func (s *userService) ChangePassword(userID uuid.UUID, req request.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(req.CurrentPassword, user.PasswordHash) {
		return errors.New("PROFILE_INVALID_CURRENT_PASSWORD")
	}

	newHashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = newHashedPassword

	return s.userRepo.Update(user)
}

func (s *userService) RegenerateAPIKey(userID uuid.UUID) (string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", err
	}

	newAPIKey, err := utils.GenerateAPIKey()
	if err != nil {
		return "", err
	}
	user.APIKey = newAPIKey

	err = s.userRepo.Update(user)
	return newAPIKey, err
}
