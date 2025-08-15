package handlers

import (
	"net/http"
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/request"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/response"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
	cfg         configs.Config
}

func NewAuthHandler(authService services.AuthService, cfg configs.Config) *AuthHandler {
	return &AuthHandler{authService: authService, cfg: cfg}
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

// Login godoc
// @Summary Log in a user
// @Description Authenticates a user and returns an access token.
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param   credentials body request.LoginRequest true "User Login Credentials"
// @Success 200 {object} response.LoginSuccessResponse "Login successful"
// @Failure 400 {object} response.APIErrorResponse "Validation error"
// @Failure 401 {object} response.APIErrorResponse "Invalid credentials"
// @Failure 500 {object} response.APIErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	loginResult, err := h.authService.Login(req)
	if err != nil {
		if err.Error() == "AUTH_INVALID_CREDENTIALS" {
			response.SendError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password", nil)
			return
		}
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to login user", nil)
		return
	}

	c.JSON(http.StatusOK, response.LoginSuccessResponse{
		Success: true,
		Message: "Login successful",
		Data: response.LoginResponse{
			AccessToken:  loginResult.AccessToken,
			RefreshToken: loginResult.RefreshToken,
			User:         response.ToUserResponse(loginResult.User),
		},
		Timestamp: time.Now().UTC(),
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generates a new access token using a refresh token.
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param   token body request.RefreshTokenRequest true "Refresh Token"
// @Success 200 {object} response.RefreshTokenSuccessResponse "Token refreshed successfully"
// @Failure 400 {object} response.APIErrorResponse "Validation error"
// @Failure 401 {object} response.APIErrorResponse "Invalid refresh token"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	newAccessToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		response.SendError(c, http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", "Invalid or expired refresh token", nil)
		return
	}

	// Ambil durasi dari config untuk ditampilkan di response
	expiresIn, _ := time.ParseDuration(h.cfg.JWT.ExpiresIn)

	c.JSON(http.StatusOK, response.RefreshTokenSuccessResponse{
		Success: true,
		Data: response.RefreshTokenResponse{
			AccessToken: newAccessToken,
			ExpiresIn:   int64(expiresIn.Seconds()),
			TokenType:   "Bearer",
		},
		Timestamp: time.Now().UTC(),
	})
}
