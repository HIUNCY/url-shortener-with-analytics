package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/HIUNCY/url-shortener-with-analytics/internal/dto/response"
	"github.com/HIUNCY/url-shortener-with-analytics/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QRCodeHandler struct {
	qrCodeService services.QRCodeService
}

func NewQRCodeHandler(qrCodeService services.QRCodeService) *QRCodeHandler {
	return &QRCodeHandler{qrCodeService: qrCodeService}
}

// GetQRCode godoc
// @Summary Get QR Code Info
// @Description Retrieves QR code as a base64 string and other info.
// @Tags QR Codes
// @Security BearerAuth
// @Produce  json
// @Param    url_id path string true "URL ID" format(uuid)
// @Param    size query int false "QR code size in pixels" default(256)
// @Success 200 {object} response.QRCodeSuccessResponse
// @Failure 403 {object} response.APIErrorResponse "Forbidden"
// @Failure 404 {object} response.APIErrorResponse "URL not found"
// @Router /urls/{url_id}/qr [get]
func (h *QRCodeHandler) GetQRCode(c *gin.Context) {
	urlID, _ := uuid.Parse(c.Param("urlID"))
	userID := c.MustGet("userID").(uuid.UUID)
	size, _ := strconv.Atoi(c.DefaultQuery("size", "256"))

	_, qrCode, err := h.qrCodeService.GetQRCodeInfo(urlID, userID, size, "png")
	if err != nil {
		// Handle not found, forbidden errors...
		response.SendError(c, http.StatusNotFound, "NOT_FOUND", "Cannot generate QR code for the URL", nil)
		return
	}

	downloadURL := fmt.Sprintf("/api/v1/urls/%s/qr/download?size=%d", urlID, size)
	c.JSON(http.StatusOK, response.QRCodeSuccessResponse{
		Success: true,
		Data: response.QRCodeResponse{
			QRCode:      qrCode,
			Format:      "png",
			Size:        size,
			DownloadURL: downloadURL,
		},
		Timestamp: time.Now().UTC(),
	})
}

// DownloadQRCode godoc
// @Summary Download QR Code
// @Description Downloads the QR code image file.
// @Tags QR Codes
// @Security BearerAuth
// @Produce  image/png
// @Param    url_id path string true "URL ID" format(uuid)
// @Param    size query int false "QR code size in pixels" default(256)
// @Success 200 {file} binary "QR Code Image"
// @Failure 403 {object} response.APIErrorResponse "Forbidden"
// @Failure 404 {object} response.APIErrorResponse "URL not found"
// @Router /urls/{url_id}/qr/download [get]
func (h *QRCodeHandler) DownloadQRCode(c *gin.Context) {
	urlID, _ := uuid.Parse(c.Param("urlID"))
	userID := c.MustGet("userID").(uuid.UUID)
	size, _ := strconv.Atoi(c.DefaultQuery("size", "256"))

	url, qrCodeBytes, err := h.qrCodeService.GetQRCodeForDownload(urlID, userID, size)
	if err != nil {
		response.SendError(c, http.StatusNotFound, "NOT_FOUND", "Cannot generate QR code for the URL", nil)
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s_qrcode.png\"", url.ShortCode))
	c.Data(http.StatusOK, "image/png", qrCodeBytes)
}
