package request

type UnlockURLRequest struct {
	Password string `json:"password" binding:"required"`
}
