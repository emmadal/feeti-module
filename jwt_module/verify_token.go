package jwt_modules

import (
	"errors"

	jwt "github.com/golang-jwt/jwt/v5"
)

// VerifyToken verify the given token to get its payload.
func VerifyToken(tokenString string, secretKey []byte) (int64, error) {
	// signature verification
	parsedToken, err := jwt.Parse(tokenString, func(tokenString *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return 0, errors.New("Invalid token")
	}

	// check if the token is valid
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		userID := claims["userID"].(float64) // convert claims to float64
		return int64(userID), nil
	} else {
		return 0, errors.New("Unable to handle token")
	}
}
