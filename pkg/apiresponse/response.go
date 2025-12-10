package apiresponse

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/hotel_mgmt_system/pkg/apperror"
)

type StandardResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

type ErrorInfo struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, StandardResponse{
		Success:   true,
		Message:   msg,
		Data:      data,
		Timestamp: nowUTC(),
		RequestID: getRequestID(c),
	})
}

func Created(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusCreated, StandardResponse{
		Success:   true,
		Message:   msg,
		Data:      data,
		Timestamp: nowUTC(),
		RequestID: getRequestID(c),
	})
}

func Deleted(c *gin.Context, msg string) {
	c.JSON(http.StatusNoContent, StandardResponse{
		Success:   true,
		Message:   msg,
		Timestamp: nowUTC(),
		RequestID: getRequestID(c),
	})
}

func Error(c *gin.Context, err error, details map[string]string) {
	var appErr *apperror.AppError

	// If it's an AppError, return structured known error
	if errors.As(err, &appErr) {
		c.JSON(appErr.HTTPStatusCode, StandardResponse{
			Success: false,
			Message: "Request failed",
			Error: &ErrorInfo{
				Code:    appErr.Code,
				Message: appErr.PublicMsg,
				Details: details,
			},
			Timestamp: nowUTC(),
			RequestID: getRequestID(c),
		})
		return
	}

	// Unknown / internal errors
	c.JSON(http.StatusInternalServerError, StandardResponse{
		Success: false,
		Message: "Request failed",
		Error: &ErrorInfo{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Something went wrong. Please try again later.",
		},
		Timestamp: nowUTC(),
		RequestID: getRequestID(c),
	})
}

func nowUTC() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func getRequestID(c *gin.Context) string {
	rid := c.Request.Header.Get("X-Request-ID")
	if rid == "" {
		return ""
	}
	return rid
}
