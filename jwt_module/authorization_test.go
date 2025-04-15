package jwt_module

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthAuthorization(t *testing.T) {
	secretKey := []byte("my_secret_key")

	// Create a valid token for testing
	validToken, err := GenerateToken(1, secretKey)
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
			expectedStatus: 403,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new gin context
			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(AuthAuthorization(secretKey))
			r.GET("/test", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "success"})
			})

			// Create a new HTTP request
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			req.AddCookie(&http.Cookie{
				Name:     "token",
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
	secretKey := []byte("my_secret_key")
	validToken, _ := GenerateToken(1, secretKey)

	// Create a new gin context
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthAuthorization(secretKey))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	for b.Loop() {
		// Create a new HTTP request
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{
			Name:     "token",
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
