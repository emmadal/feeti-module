package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID int64 `json:"userID"`
	jwt.RegisteredClaims
}

// GenerateToken generate a valid jwt token for 30 minutes
func GenerateToken(userID int64, secretKey []byte) (string, error) {
	// check if the secret key and userID are valid
	if userID <= 0 || len(secretKey) == 0 {
		return "", fmt.Errorf("invalid user id")
	}
	// create a new token with the given userID
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Minute)),
		},
	})
	// sign the token with the secret key
	return token.SignedString(secretKey)
}
