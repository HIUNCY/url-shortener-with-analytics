package postgres

import (
	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru dari userRepository.
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Store(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

// Implementasi method lain dari interface bisa ditambahkan di sini nanti.
func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}
func (r *userRepository) FindByAPIKey(apiKey string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("api_key = ?", apiKey).First(&user).Error
	return &user, err
}
func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}
