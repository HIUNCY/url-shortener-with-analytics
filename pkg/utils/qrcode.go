package utils

import (
	"encoding/base64"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCodeBase64(text string, size int) (string, error) {
	png, err := qrcode.Encode(text, qrcode.Medium, size)
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png), nil
}

func GenerateQRCodeBytes(text string, size int) ([]byte, error) {
	return qrcode.Encode(text, qrcode.Medium, size)
}
