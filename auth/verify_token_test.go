package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyToken(t *testing.T) {
	token, err := GenerateToken(1, []byte("my_secret_key"))
	if assert.NoError(t, err) {
		assert.NotEmpty(t, token)
	}

	_, err = VerifyToken(token, []byte("my_secret_key"))
	assert.NoError(t, err)

}

func BenchmarkVerifyToken(b *testing.B) {
	secretKey := []byte("my_secret_key") // allocated once
	token, err := GenerateToken(1, secretKey)
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
