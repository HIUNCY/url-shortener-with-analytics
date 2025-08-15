package response

import "time"

// LoginResponse adalah DTO untuk payload data pada respons login sukses.
type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

// LoginSuccessResponse adalah struktur top-level untuk respons sukses login.
type LoginSuccessResponse struct {
	Success   bool          `json:"success" example:"true"`
	Message   string        `json:"message" example:"Login successful"`
	Data      LoginResponse `json:"data"`
	Timestamp time.Time     `json:"timestamp"`
}
