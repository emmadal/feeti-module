package auth

import (
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var userID = uuid.New()
var secretKey = []byte("my_secret_key")

func TestAuthAuthorization(t *testing.T) {
	// Create a valid token for testing
	validToken, err := GenerateToken(userID, secretKey)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{
			name:           "Valid token",
			token:          validToken,
			expectedStatus: 200,
		},
		{
			name:           "Empty token",
			token:          "",
			expectedStatus: 401,
		},
		{
			name:           "Invalid token",
			token:          "invalid_token",
			expectedStatus: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new gin context
			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(AuthGin(secretKey))
			r.GET("/test", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "success"})
			})

			// Create a new HTTP request
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			req.AddCookie(&http.Cookie{
				Name:     "ftk",
				Value:    tt.token,
				Domain:   "/",
				HttpOnly: true,
				Secure:   false,
				Path:     "/",
			})

			// Create a response recorder
			w := httptest.NewRecorder()

			// Serve the request
			r.ServeHTTP(w, req)

			// Assert the response status code
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func BenchmarkAuthAuthorization(b *testing.B) {
	validToken, _ := GenerateToken(userID, secretKey)

	// Create a new gin context
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthGin(secretKey))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	for b.Loop() {
		// Create a new HTTP request
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{
			Name:     "ftk",
			Value:    validToken,
			Domain:   "/",
			HttpOnly: true,
			Secure:   false,
			Path:     "/",
		})

		// Create a response recorder
		w := httptest.NewRecorder()

		// Serve the request
		r.ServeHTTP(w, req)
	}
}

func TestGetUserIDFromGin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthGin(secretKey))
	r.GET("/test", func(c *gin.Context) {
		contextUserID := GetUserIDFromGin(c)
		assert.Equal(t, userID, contextUserID)
		c.JSON(200, gin.H{"userID": userID})
	})
}
