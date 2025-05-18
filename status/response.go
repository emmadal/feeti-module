package helpers

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

// HandleError is a helper function to handle an error
func HandleError(c *gin.Context, status int, message string, err error) {
	logger.Error(message)
	c.SecureJSON(
		status, gin.H{
			"message": message,
			"success": false,
		},
	)
}

// HandleSuccess is a helper function to handle a success
func HandleSuccess(c *gin.Context, message string) {
	c.SecureJSON(
		http.StatusOK, gin.H{
			"message": message,
			"success": true,
		},
	)
}

// HandleSuccessData is a helper function to handle a success and data
func HandleSuccessData(c *gin.Context, message string, data any) {
	c.SecureJSON(
		http.StatusOK, gin.H{
			"message": message,
			"success": true,
			"data":    data,
		},
	)
}
