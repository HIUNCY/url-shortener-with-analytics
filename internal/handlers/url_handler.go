package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/configs"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/domain"
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

// GetUserURLs godoc
// @Summary Get user's URLs
// @Description Retrieves a paginated list of URLs for the authenticated user.
// @Tags URLs
// @Security BearerAuth
// @Security ApiKeyAuth
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search query for title or original URL"
// @Param sort query string false "Sort by field (created_at, click_count, title)" Enums(created_at, click_count, title)
// @Param order query string false "Sort order (asc, desc)" Enums(asc, desc)
// @Success 200 {object} response.URLListSuccessResponse "List of URLs retrieved successfully"
// @Failure 401 {object} response.APIErrorResponse "Unauthorized"
// @Router /urls [get]
func (h *URLHandler) GetUserURLs(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	// Parsing query parameters dengan nilai default
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	options := &domain.FindAllOptions{
		Search: c.Query("search"),
		SortBy: c.DefaultQuery("sort", "created_at"),
		Order:  c.DefaultQuery("order", "desc"),
		Limit:  limit,
		Offset: offset,
	}

	result, err := h.urlService.GetUserURLs(userID, options)
	if err != nil {
		response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to retrieve URLs", nil)
		return
	}

	// Mapping dari domain.URL ke DTO
	urlResponses := make([]response.URLListItemResponse, len(result.URLs))
	for i, url := range result.URLs {
		shortURLString := fmt.Sprintf("%s/%s", h.cfg.Server.BaseURL, url.ShortCode)
		urlResponses[i] = response.URLListItemResponse{
			ID:          url.ID,
			OriginalURL: url.OriginalURL,
			ShortCode:   url.ShortCode,
			ShortURL:    shortURLString,
			Title:       url.Title,
			ClickCount:  url.ClickCount,
			IsActive:    url.IsActive,
			ExpiresAt:   url.ExpiresAt,
			CreatedAt:   url.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response.URLListSuccessResponse{
		Success: true,
		Data: response.URLListResponse{
			URLs:       urlResponses,
			Pagination: result.Pagination,
		},
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

// UpdateURL godoc
// @Summary Update a URL
// @Description Updates the properties of a specific short URL.
// @Tags URLs
// @Security BearerAuth
// @Security ApiKeyAuth
// @Accept   json
// @Produce  json
// @Param    url_id path string true "URL ID" format(uuid)
// @Param    url body request.UpdateURLRequest true "URL Update Information"
// @Success 200 {object} response.URLDetailsSuccessResponse "URL updated successfully"
// @Failure 400 {object} response.APIErrorResponse "Validation error"
// @Failure 403 {object} response.APIErrorResponse "Forbidden"
// @Failure 404 {object} response.APIErrorResponse "URL not found"
// @Router /urls/{url_id} [put]
func (h *URLHandler) UpdateURL(c *gin.Context) {
	urlID, err := uuid.Parse(c.Param("urlID"))
	if err != nil {
		response.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid URL ID format", nil)
		return
	}

	var req request.UpdateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	updatedURL, err := h.urlService.UpdateURL(urlID, userID, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.SendError(c, http.StatusNotFound, "NOT_FOUND", "URL not found", nil)
			return
		}
		if err.Error() == "URL_FORBIDDEN" {
			response.SendError(c, http.StatusForbidden, "FORBIDDEN", "You do not have permission to view this URL", nil)
			return
		}
		response.SendError(c, http.StatusInternalServerError, "UPDATE_FAILED", "Failed to update URL", nil)
		return
	}

	shortURLString := fmt.Sprintf("%s/%s", h.cfg.Server.BaseURL, updatedURL.ShortCode)
	c.JSON(http.StatusOK, response.URLDetailsSuccessResponse{
		Success:   true,
		Data:      response.ToURLDetailsResponse(updatedURL, shortURLString),
		Timestamp: time.Now().UTC(),
	})
}

// DeleteURL godoc
// @Summary Delete a URL
// @Description Deletes a specific short URL (soft delete).
// @Tags URLs
// @Security BearerAuth
// @Security ApiKeyAuth
// @Produce  json
// @Param    url_id path string true "URL ID" format(uuid)
// @Success 200 {object} response.SuccessMessageResponse "URL deleted successfully"
// @Failure 403 {object} response.APIErrorResponse "Forbidden"
// @Failure 404 {object} response.APIErrorResponse "URL not found"
// @Router /urls/{url_id} [delete]
func (h *URLHandler) DeleteURL(c *gin.Context) {
	urlID, err := uuid.Parse(c.Param("urlID"))
	if err != nil {
		response.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid URL ID format", nil)
		return
	}

	userID := c.MustGet("userID").(uuid.UUID)
	if err := h.urlService.DeleteURL(urlID, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.SendError(c, http.StatusNotFound, "NOT_FOUND", "URL not found", nil)
			return
		}
		if err.Error() == "URL_FORBIDDEN" {
			response.SendError(c, http.StatusForbidden, "FORBIDDEN", "You do not have permission to view this URL", nil)
			return
		}
		response.SendError(c, http.StatusInternalServerError, "DELETE_FAILED", "Failed to delete URL", nil)
		return
	}

	c.JSON(http.StatusOK, response.SuccessMessageResponse{
		Success:   true,
		Message:   "URL deleted successfully",
		Timestamp: time.Now().UTC(),
	})
}
