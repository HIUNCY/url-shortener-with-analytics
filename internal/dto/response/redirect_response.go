package response

import "time"

// UnlockURLResponse adalah DTO untuk payload data pada respons unlock sukses.
type UnlockURLResponse struct {
	RedirectURL string `json:"redirect_url"`
	AccessToken string `json:"access_token"`
}

// UnlockURLSuccessResponse adalah wrapper untuk Swagger.
type UnlockURLSuccessResponse struct {
	Success   bool              `json:"success" example:"true"`
	Data      UnlockURLResponse `json:"data"`
	Timestamp time.Time         `json:"timestamp"`
}
