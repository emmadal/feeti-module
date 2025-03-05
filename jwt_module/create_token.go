package jwt_modules

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateToken generate a valid jwt token for 30 minutes
func GenerateToken(userID int64, secretKey []byte) (string, error) {
	// check if the secret key and userID are valid
	if userID <= 0 || len(secretKey) == 0 {
		return "", errors.New("invalid secret key or userID")
	}
	// create a new token with the given userID
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"iat":    now.Unix(),
		"exp":    now.Add(30 * time.Minute).Unix(),
		"jti":    uuid.New().String(),
	})
	// sign the token with the secret key
	return token.SignedString(secretKey)
}
