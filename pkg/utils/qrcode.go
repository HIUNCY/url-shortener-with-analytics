package utils

import (
	"encoding/base64"

	"github.com/skip2/go-qrcode"
)

// GenerateQRCodeBase64 membuat QR code dari sebuah teks dan mengembalikannya sebagai string base64 yang siap disisipkan di HTML/JSON.
func GenerateQRCodeBase64(text string, size int) (string, error) {
	png, err := qrcode.Encode(text, qrcode.Medium, size)
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png), nil
}
