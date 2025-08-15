package response

import (
	"time"

	"github.com/gin-gonic/gin"
)

// APIResponse adalah struktur dasar untuk semua respons sukses.
// Menggunakan generic type T untuk data payload.
type APIResponse[T any] struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Data      T         `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// ErrorDetail mendeskripsikan error validasi untuk satu field.
type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// APIErrorResponse adalah struktur untuk semua respons error.
type APIErrorResponse struct {
	Success   bool         `json:"success"`
	Error     ErrorPayload `json:"error"`
	Timestamp time.Time    `json:"timestamp"`
	RequestID string       `json:"request_id,omitempty"`
}

// ErrorPayload berisi detail dari error.
type ErrorPayload struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type SuccessMessageResponse struct {
	Success   bool      `json:"success" example:"true"`
	Message   string    `json:"message" example:"Operation successful"`
	Timestamp time.Time `json:"timestamp"`
}

// SendSuccess adalah helper untuk mengirim respons sukses.
func SendSuccess[T any](c *gin.Context, statusCode int, message string, data T) {
	c.JSON(statusCode, APIResponse[T]{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	})
}

// SendError adalah helper untuk mengirim respons error.
func SendError(c *gin.Context, statusCode int, code, message string, details []ErrorDetail) {
	c.AbortWithStatusJSON(statusCode, APIErrorResponse{
		Success: false,
		Error: ErrorPayload{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now().UTC(),
	})
}
