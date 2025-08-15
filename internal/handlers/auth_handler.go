package handlers

import (
	"net/http"
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/request"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/response"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account with the provided details.
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param   user body request.RegisterRequest true "User Registration Info"
// @Success 201 {object} response.RegisterSuccessResponse "User registered successfully" // <-- DIUBAH DI SINI
// @Failure 400 {object} response.APIErrorResponse "Validation error"
// @Failure 409 {object} response.APIErrorResponse "Email already exists"
// @Failure 500 {object} response.APIErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Kita akan buat error detail nanti. Untuk sekarang, nil saja.
		response.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	newUser, err := h.authService.Register(req)
	if err != nil {
		if err.Error() == "AUTH_EMAIL_ALREADY_EXISTS" {
			response.SendError(c, http.StatusConflict, "EMAIL_CONFLICT", "User with this email already exists", nil)
			return
		}
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to register user", nil)
		return
	}

	// Kirim respons sukses menggunakan struct baru.
	// Helper generik kita hapus untuk sementara agar lebih jelas.
	c.JSON(http.StatusCreated, response.RegisterSuccessResponse{
		Success:   true,
		Message:   "User registered successfully",
		Data:      response.ToUserResponse(newUser),
		Timestamp: time.Now().UTC(),
	})
}
