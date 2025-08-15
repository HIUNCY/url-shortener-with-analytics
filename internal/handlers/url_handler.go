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

type URLHandler struct {
	urlService services.URLService
}

func NewURLHandler(urlService services.URLService) *URLHandler {
	return &URLHandler{urlService: urlService}
}

// CreateShortURL godoc
// @Summary Create a new short URL
// @Description Creates a new short URL for the authenticated user.
// @Tags URLs
// @Security BearerAuth
// @Security ApiKeyAuth
// @Accept   json
// @Produce  json
// @Param    url body request.CreateURLRequest true "URL Information"
// @Success 201 {object} response.CreateURLSuccessResponse "URL created successfully"
// @Failure 400 {object} response.APIErrorResponse "Validation error"
// @Failure 401 {object} response.APIErrorResponse "Unauthorized"
// @Failure 409 {object} response.APIErrorResponse "Custom alias already exists"
// @Router /urls [post]
func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var req request.CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	result, err := h.urlService.CreateShortURL(userID, req)
	if err != nil {
		if err.Error() == "URL_CUSTOM_ALIAS_EXISTS" {
			response.SendError(c, http.StatusConflict, "ALIAS_CONFLICT", "Custom alias already exists", nil)
			return
		}
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to create short URL", nil)
		return
	}

	c.JSON(http.StatusCreated, response.CreateURLSuccessResponse{
		Success:   true,
		Message:   "Short URL created successfully",
		Data:      response.ToCreateURLResponse(result.URL, result.ShortURL, result.QRCode),
		Timestamp: time.Now().UTC(),
	})
}
