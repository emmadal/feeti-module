package jwt_module

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthAuthorization is middleware for JWT authentication
func AuthAuthorization(secretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate secret key
		if len(secretKey) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		// Get the token from the cookie
		tokenCookie, err := c.Request.Cookie("ftk")
		if err != nil || tokenCookie.Value == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
			return
		}

		// Get the token value
		token := strings.TrimSpace(tokenCookie.Value)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authentication required"})
			return
		}

		// Verify the token
		userID, err := VerifyToken(token, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Authentication failed"})
			return
		}

		// Attach userID to the context
		c.Set("userID", userID)
		c.Next()
	}
}
