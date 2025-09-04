package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SetSecureCookie sets a JWT token in a cookie with secure settings
func SetSecureCookie(c *gin.Context, token string, domain string) {
	var sameSite http.SameSite
	// Create a new cookie with the token
	if domain == "localhost" {
		sameSite = http.SameSiteLaxMode
	} else {
		sameSite = http.SameSiteNoneMode
	}
	cookie := &http.Cookie{
		Name:     "ftk",
		Value:    token,
		Path:     "/",
		Domain:   domain,
		MaxAge:   int(30 * time.Minute), // Match token expiration time
		HttpOnly: true,                  // Prevent JavaScript access
		SameSite: sameSite,
	}

	// Set Secure based on HTTPS usage
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		cookie.Secure = true
	} else {
		cookie.Secure = false
	}

	// Set the cookie in the response
	http.SetCookie(c.Writer, cookie)
}

// ClearAuthCookie clears the authentication cookie
func ClearAuthCookie(c *gin.Context, domain string) {
	var sameSite http.SameSite
	// Create a new cookie with the token
	if domain == "localhost" {
		sameSite = http.SameSiteLaxMode
	} else {
		sameSite = http.SameSiteNoneMode
	}
	cookie := &http.Cookie{
		Name:     "ftk",
		Value:    "",
		Path:     "/",
		Domain:   domain,
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: sameSite,
	}
	// Set Secure based on HTTPS usage
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		cookie.Secure = true
	} else {
		cookie.Secure = false
	}
	http.SetCookie(c.Writer, cookie)
}
