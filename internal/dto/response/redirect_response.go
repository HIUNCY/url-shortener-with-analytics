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

type URLInfoResponse struct {
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	ClickCount  int       `json:"click_count"`
	CreatedAt   time.Time `json:"created_at"`
	IsSafe      bool      `json:"is_safe"`
	Domain      string    `json:"domain"`
}

type URLInfoSuccessResponse struct {
	Success   bool            `json:"success" example:"true"`
	Data      URLInfoResponse `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
}
