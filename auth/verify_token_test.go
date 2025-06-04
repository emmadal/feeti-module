package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyToken(t *testing.T) {

	token, err := GenerateToken(userID, secretKey)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, token)
	}

	_, err = VerifyToken(token, secretKey)
	assert.NoError(t, err)

}

func BenchmarkVerifyToken(b *testing.B) {
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
