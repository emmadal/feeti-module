package jwt_modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken(1, []byte("my_secret_key"))
	if err != nil {
		t.Errorf("Error generating token: %v", err)
	}
	assert.NotEmpty(t, token)
}
