package response

import "time"

type QRCodeResponse struct {
	QRCode      string `json:"qr_code"`
	Format      string `json:"format"`
	Size        int    `json:"size"`
	DownloadURL string `json:"download_url"`
}

type QRCodeSuccessResponse struct {
	Success   bool           `json:"success" example:"true"`
	Data      QRCodeResponse `json:"data"`
	Timestamp time.Time      `json:"timestamp"`
}
