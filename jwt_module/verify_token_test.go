package jwt_modules

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
