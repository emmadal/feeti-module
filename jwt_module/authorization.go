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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid secret key"})
			return
		}

		// Get the token from the request header
		tokenCookie, err := c.Request.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token not found"})
			return
		}

		token := strings.TrimSpace(tokenCookie.Value)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Empty token"})
			return
		}

		// Verify the token
		userID, err := VerifyToken(token, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Invalid token"})
			return
		}

		// Attach userID to the context
		c.Set("userID", userID)
		c.Next()
	}
}
