package jwt_module

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken(1, []byte("my_secret_key"))
	if err != nil {
		t.Errorf("Error generating token: %v", err)
	}
	assert.NotEmpty(t, token)
}

func BenchmarkGenerateToken(b *testing.B) {
	// Setup code outside the measured part
	secretKey := []byte("my_secret_key")
	userID := int64(1)

	// ReportAllocs will report memory allocations
	b.ReportAllocs()

	// Run the benchmark with b.N iterations
	for n := 0; n < b.N; n++ {
		_, err := GenerateToken(userID, secretKey)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGenerateTokenAlloc measures memory allocations
func BenchmarkGenerateTokenAlloc(b *testing.B) {
	secretKey := []byte("my_secret_key")
	userID := int64(1)

	// ReportAllocs will report memory allocations
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		_, err := GenerateToken(userID, secretKey)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCreateClaims focuses specifically on the claims creation portion
func BenchmarkCreateClaims(b *testing.B) {
	userID := int64(1)

	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		now := time.Now()
		_ = &UserClaims{
			UserID: userID,
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(now),
				ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Minute)),
			},
		}
	}
}

// BenchmarkNewWithClaims focuses on the JWT creation step
func BenchmarkNewWithClaims(b *testing.B) {
	userID := int64(1)
	var tokens []*jwt.Token

	// Pre-create claims objects to isolate the token creation part
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		now := time.Now()
		claims := &UserClaims{
			UserID: userID,
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(now),
				ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Minute)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokens = append(tokens, token) // prevent GC from collecting too early
	}

	// Use tokens to prevent compiler optimizations
	b.StopTimer()
	if len(tokens) == 0 {
		b.Fatal("tokens should not be empty")
	}
}

// BenchmarkSignedString focuses on the JWT signing step, which showed highest allocations
func BenchmarkSignedString(b *testing.B) {
	secretKey := []byte("my_secret_key")
	userID := int64(1)

	// Pre-create a token to isolate the signing part
	now := time.Now()
	claims := &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		_, err := token.SignedString(secretKey)
		if err != nil {
			b.Fatal(err)
		}
	}
}
