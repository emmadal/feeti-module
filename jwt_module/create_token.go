package jwt_modules

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken generate a valid jwt token for 10 minutes
func GenerateToken(userID int64, secretKey []byte) (string, error) {
	// check if the secret key and userID are valid
	if userID <= 0 || secretKey == nil {
		return "", errors.New("invalid secret key or userID")
	}
	// create a new token with the given userID
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(10 * time.Minute).Unix(), // expire in 10 minutes
	})
	// sign the token with the secret key
	return token.SignedString(secretKey)
}
