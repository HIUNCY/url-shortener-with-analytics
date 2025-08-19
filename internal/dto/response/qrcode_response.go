package response

import "time"

// QRCodeResponse adalah DTO untuk payload data pada respons info QR code.
type QRCodeResponse struct {
	QRCode      string `json:"qr_code"`
	Format      string `json:"format"`
	Size        int    `json:"size"`
	DownloadURL string `json:"download_url"`
}

// QRCodeSuccessResponse adalah wrapper untuk Swagger.
type QRCodeSuccessResponse struct {
	Success   bool           `json:"success" example:"true"`
	Data      QRCodeResponse `json:"data"`
	Timestamp time.Time      `json:"timestamp"`
}
