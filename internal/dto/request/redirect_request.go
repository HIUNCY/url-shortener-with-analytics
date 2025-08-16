package request

// UnlockURLRequest adalah DTO untuk request membuka URL yang terproteksi.
type UnlockURLRequest struct {
	Password string `json:"password" binding:"required"`
}
