package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/request"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/response"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type URLHandler struct {
	urlService services.URLService
	cfg        configs.Config
}

func NewURLHandler(urlService services.URLService, cfg configs.Config) *URLHandler {
	return &URLHandler{urlService: urlService, cfg: cfg}
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

// GetURLDetails godoc
// @Summary Get URL details
// @Description Retrieves the details of a specific short URL owned by the user.
// @Tags URLs
// @Security BearerAuth
// @Security ApiKeyAuth
// @Produce  json
// @Param    url_id path string true "URL ID" format(uuid)
// @Success 200 {object} response.URLDetailsSuccessResponse "URL details retrieved successfully"
// @Failure 401 {object} response.APIErrorResponse "Unauthorized"
// @Failure 403 {object} response.APIErrorResponse "Forbidden"
// @Failure 404 {object} response.APIErrorResponse "URL not found"
// @Router /urls/{url_id} [get]
func (h *URLHandler) GetURLDetails(c *gin.Context) {
	// 1. Ambil parameter dari URL
	urlID, err := uuid.Parse(c.Param("urlID"))
	if err != nil {
		response.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid URL ID format", nil)
		return
	}

	// 2. Ambil userID dari context (di-set oleh middleware)
	userID := c.MustGet("userID").(uuid.UUID)

	// 3. Panggil service
	url, err := h.urlService.GetURLDetails(urlID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.SendError(c, http.StatusNotFound, "NOT_FOUND", "URL not found", nil)
			return
		}
		if err.Error() == "URL_FORBIDDEN" {
			response.SendError(c, http.StatusForbidden, "FORBIDDEN", "You do not have permission to view this URL", nil)
			return
		}
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to retrieve URL", nil)
		return
	}

	shortURLString := fmt.Sprintf("%s/%s", h.cfg.Server.BaseURL, url.ShortCode)

	// 4. Kirim respons sukses
	c.JSON(http.StatusOK, response.URLDetailsSuccessResponse{
		Success:   true,
		Data:      response.ToURLDetailsResponse(url, shortURLString),
		Timestamp: time.Now().UTC(),
	})
}
