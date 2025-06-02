package auth

import (
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyToken(t *testing.T) {
	userID := uuid.New()
	token, err := GenerateToken(userID, []byte("my_secret_key"))
	if assert.NoError(t, err) {
		assert.NotEmpty(t, token)
	}

	_, err = VerifyToken(token, []byte("my_secret_key"))
	assert.NoError(t, err)

}

func BenchmarkVerifyToken(b *testing.B) {
	secretKey := []byte("my_secret_key") // allocated once
	userID := uuid.New()
	token, err := GenerateToken(userID, secretKey)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := VerifyToken(token, secretKey)
		if err != nil {
			b.Fatal(err)
		}
	}
}
