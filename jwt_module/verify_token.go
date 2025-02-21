package jwt_modules

import (
	"errors"

	jwt "github.com/golang-jwt/jwt/v5"
)

// VerifyToken verify the given token to get its payload.
func VerifyToken(tokenString string, secretKey []byte) (int64, error) {

	// verify the signature of the token
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Invalid signature")
		}
		return secretKey, nil
	})

	// check if there is an error
	if err != nil {
		return 0, errors.New("Couldn't handle this token")
	}

	// check the validity of token
	if tokenIsValid := parsedToken.Valid; !tokenIsValid {
		return 0, errors.New("Invalid or Expired token")
	}

	// extract claims data in the token
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Invalid claims")
	}
	// convert userID to uint
	userID := claims["userID"].(float64)
	return int64(userID), nil
}
