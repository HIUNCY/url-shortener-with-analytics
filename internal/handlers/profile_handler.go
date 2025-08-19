package handlers

import (
	"net/http"
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/request"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/response"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	userService services.UserService
}

func NewProfileHandler(userService services.UserService) *ProfileHandler {
	return &ProfileHandler{userService: userService}
}

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieves the profile of the currently logged-in user.
// @Tags Profile
// @Security BearerAuth
// @Produce  json
// @Success 200 {object} response.ProfileSuccessResponse
// @Failure 401 {object} response.APIErrorResponse "Unauthorized"
// @Router /profile [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	user, err := h.userService.GetProfile(userID)
	if err != nil {
		response.SendError(c, http.StatusNotFound, "NOT_FOUND", "User profile not found", nil)
		return
	}
	response.SendSuccess(c, http.StatusOK, "Profile retrieved successfully", response.ToUserResponse(user))
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Updates the first and last name of the currently logged-in user.
// @Tags Profile
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    profile body request.UpdateProfileRequest true "Profile Information"
// @Success 200 {object} response.ProfileSuccessResponse
// @Failure 400 {object} response.APIErrorResponse "Validation error"
// @Failure 401 {object} response.APIErrorResponse "Unauthorized"
// @Router /profile [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	userID := c.MustGet("userID").(uuid.UUID)
	updatedUser, err := h.userService.UpdateProfile(userID, req)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, "UPDATE_FAILED", "Failed to update profile", nil)
		return
	}
	response.SendSuccess(c, http.StatusOK, "Profile updated successfully", response.ToUserResponse(updatedUser))
}

// ChangePassword godoc
// @Summary Change user password
// @Description Updates the password of the currently logged-in user.
// @Tags Profile
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    passwords body request.ChangePasswordRequest true "Password Change Info"
// @Success 200 {object} response.SuccessMessageResponse
// @Failure 400 {object} response.APIErrorResponse "Validation error"
// @Failure 401 {object} response.APIErrorResponse "Invalid current password"
// @Router /profile/password [put]
func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	userID := c.MustGet("userID").(uuid.UUID)
	err := h.userService.ChangePassword(userID, req)
	if err != nil {
		if err.Error() == "PROFILE_INVALID_CURRENT_PASSWORD" {
			response.SendError(c, http.StatusUnauthorized, "INVALID_PASSWORD", "Current password is incorrect", nil)
			return
		}
		response.SendError(c, http.StatusInternalServerError, "UPDATE_FAILED", "Failed to change password", nil)
		return
	}

	c.JSON(http.StatusOK, response.SuccessMessageResponse{
		Success:   true,
		Message:   "Password changed successfully",
		Timestamp: time.Now().UTC(),
	})
}

// RegenerateAPIKey godoc
// @Summary Regenerate API Key
// @Description Generates a new API key for the currently logged-in user.
// @Tags Profile
// @Security BearerAuth
// @Produce  json
// @Success 200 {object} response.APIKeySuccessResponse
// @Failure 401 {object} response.APIErrorResponse "Unauthorized"
// @Router /profile/api-key/regenerate [post]
func (h *ProfileHandler) RegenerateAPIKey(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	newAPIKey, err := h.userService.RegenerateAPIKey(userID)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, "GENERATION_FAILED", "Failed to regenerate API key", nil)
		return
	}

	c.JSON(http.StatusOK, response.APIKeySuccessResponse{
		Success:   true,
		Message:   "API key regenerated successfully",
		Data:      response.APIKeyResponse{APIKey: newAPIKey, GeneratedAt: time.Now().UTC()},
		Timestamp: time.Now().UTC(),
	})
}
