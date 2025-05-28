package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

const defaultTimeout = 5 * time.Second

// Timeout middleware is used to handle HTTP requests that take too long to
// process.
// It returns a 408 status indicating the request timed out.
// By default, the timeout duration is set to 5 seconds.
func Timeout(duration time.Duration) gin.HandlerFunc {
	if duration <= 0 {
		duration = defaultTimeout
	}
	return timeout.New(
		timeout.WithTimeout(duration),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			c.SecureJSON(http.StatusRequestTimeout, gin.H{
				"message": "Request timed out",
			})
		}),
	)
}
