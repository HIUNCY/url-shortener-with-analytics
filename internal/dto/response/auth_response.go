package response

import "time"

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type LoginSuccessResponse struct {
	Success   bool          `json:"success" example:"true"`
	Message   string        `json:"message" example:"Login successful"`
	Data      LoginResponse `json:"data"`
	Timestamp time.Time     `json:"timestamp"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type" example:"Bearer"`
}

type RefreshTokenSuccessResponse struct {
	Success   bool                 `json:"success" example:"true"`
	Data      RefreshTokenResponse `json:"data"`
	Timestamp time.Time            `json:"timestamp"`
}
