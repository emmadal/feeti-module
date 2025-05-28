package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// SetSecureCookie sets a JWT token in a cookie with secure settings
func SetSecureCookie(c *gin.Context, token string, domain string, isProduction bool) {
	// Create a new cookie with the token
	cookie := &http.Cookie{
		Name:     "ftk",
		Value:    token,
		Path:     "/",
		Domain:   domain,
		MaxAge:   int(30 * time.Minute), // Match token expiration time
		HttpOnly: true,                  // Prevent JavaScript access
		SameSite: http.SameSiteLaxMode,
	}

	// Set Secure based on HTTPS usage
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		cookie.Secure = true
	}

	// Set the cookie in the response
	http.SetCookie(c.Writer, cookie)
}

// ClearAuthCookie clears the authentication cookie
func ClearAuthCookie(c *gin.Context, domain string) {
	cookie := &http.Cookie{
		Name:     "ftk",
		Value:    "",
		Path:     "/",
		Domain:   domain,
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // Changed from SameSiteStrictMode to SameSiteLaxMode to allow cross-origin requests
	}
	http.SetCookie(c.Writer, cookie)
}
