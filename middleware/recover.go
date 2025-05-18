package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

// Recover recovers from panics and returns a 500 Internal Server Error response.
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%v", err)
				logger.Error(message)
				// Try to write header only if not already written
				if !c.Writer.Written() {
					c.AbortWithStatusJSON(
						http.StatusInternalServerError, gin.H{
							"success": false,
							"message": "Internal server error",
						},
					)
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		c.Next()
	}
}
