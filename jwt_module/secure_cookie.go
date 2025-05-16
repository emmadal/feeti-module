package jwt_module

import (
	"net/http"
	"time"
)

// SetSecureCookie sets a JWT token in a cookie with secure settings
func SetSecureCookie(w http.ResponseWriter, token string, domain string, isProduction bool) {
	// Create a new cookie with the token
	cookie := &http.Cookie{
		Name:     "ftk",
		Value:    token,
		Path:     "/",
		Domain:   domain,
		MaxAge:   int(30 * time.Minute.Seconds()), // Match token expiration time
		HttpOnly: true,                            // Prevent JavaScript access
		SameSite: http.SameSiteLaxMode,            // Protect against CSRF
	}

	// Only set a Secure flag in production, or when HTTPS is used,
	// This ensures cookies are only sent over HTTPS
	if isProduction {
		cookie.Secure = true
	}

	// Set the cookie in the response
	http.SetCookie(w, cookie)
}

// ClearAuthCookie clears the authentication cookie
func ClearAuthCookie(w http.ResponseWriter, domain string) {
	cookie := &http.Cookie{
		Name:     "ftk",
		Value:    "",
		Path:     "/",
		Domain:   domain,
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}
