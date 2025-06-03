package auth

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// UserIDKey is the context key for the user ID
type UserIDKey string

// ContextUserID is the key used to store/retrieve the user ID from context
const ContextUserID UserIDKey = "userID"

// AuthMiddleware is a middleware that checks if the user is authenticated for standard HTTP servers
func AuthMiddleware(secretKey []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate secret key
			if len(secretKey) == 0 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get the token from the cookie
			tokenCookie, err := r.Cookie("ftk")
			if err != nil || tokenCookie.Value == "" {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			// Verify the tokenCookie
			userID, err := VerifyToken(tokenCookie.Value, secretKey)
			if err != nil {
				http.Error(w, "Authentication failed", http.StatusForbidden)
				return
			}

			// Attach userID to the context and call the next handler
			ctx := context.WithValue(r.Context(), ContextUserID, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves the user ID from the request context for standard HTTP
func GetUserID(r *http.Request) uuid.UUID {
	userID, ok := r.Context().Value(ContextUserID).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return userID
}

// AuthGin is a middleware that checks if the user is authenticated for a Gin framework
func AuthGin(secretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate secret key
		if len(secretKey) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authentication failed"})
			return
		}

		// Get the token from the cookie
		tokenCookie, err := c.Request.Cookie("ftk")
		if err != nil || tokenCookie.Value == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authentication required"})
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
