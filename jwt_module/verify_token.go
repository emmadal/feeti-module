package jwt_module

import (
	"fmt"

	jwt "github.com/golang-jwt/jwt/v5"
)

// UserClaims defines a struct for JWT claims to avoid using MapClaims
type UserClaims struct {
	UserID int64 `json:"userID"`
	jwt.RegisteredClaims
}

// Global parser instance to be reused
var jwtParser = jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

// VerifyToken verify the given token to get its payload.
func VerifyToken(tokenString string, secretKey []byte) (int64, error) {
	// Create a claims instance to unmarshal into
	claims := &UserClaims{}

	// Use ParseWithClaims with the defined struct to reduce allocations
	token, err := jwtParser.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return secretKey, nil
		},
	)

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Access userID directly from the struct
	return claims.UserID, nil
}
