package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// AuthGin is a middleware that checks if the user is authenticated for a Gin framework
func AuthGin(secretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate secret key
		if len(secretKey) == 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}

		// Get the token from the cookie
		tokenCookie, err := c.Request.Cookie("ftk")
		fmt.Println("tokenCookie: ", tokenCookie)
		fmt.Println("errCookies: ", err)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		// Verify if the token is empty
		if tokenCookie.Value == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
			return
		}

		// Verify the token
		userID, err := VerifyToken(tokenCookie.Value, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authentication failed"})
			return
		}

		// Attach userID to the gin context
		c.Set("userID", userID)
		c.Next()
	}
}

// GetUserIDFromGin retrieves the user ID from the Gin context
func GetUserIDFromGin(c *gin.Context) uuid.UUID {
	userID, exists := c.Get("userID")
	if !exists {
		return uuid.Nil
	}
	return uuid.MustParse(userID.(string))
}
