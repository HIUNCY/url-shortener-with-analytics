package response

import (
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	PlanType  string    `json:"plan_type"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterSuccessResponse struct {
	Success   bool         `json:"success" example:"true"`
	Message   string       `json:"message" example:"User registered successfully"`
	Data      UserResponse `json:"data"`
	Timestamp time.Time    `json:"timestamp"`
}

func ToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		PlanType:  user.PlanType,
		APIKey:    user.APIKey,
		CreatedAt: user.CreatedAt,
	}
}
