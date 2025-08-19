package response

import (
	"time"

	"github.com/gin-gonic/gin"
)

type APIResponse[T any] struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Data      T         `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type APIErrorResponse struct {
	Success   bool         `json:"success"`
	Error     ErrorPayload `json:"error"`
	Timestamp time.Time    `json:"timestamp"`
	RequestID string       `json:"request_id,omitempty"`
}

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

func SendSuccess[T any](c *gin.Context, statusCode int, message string, data T) {
	c.JSON(statusCode, APIResponse[T]{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	})
}

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
