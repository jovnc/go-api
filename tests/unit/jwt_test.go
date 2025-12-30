package unit

import (
	"testing"

	"go_api/internal/util"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	secret := []byte("test-secret-key")

	t.Run("should generate token successfully", func(t *testing.T) {
		userID := uint(1)
		username := "testuser"

		token, err := util.GenerateToken(userID, username, secret)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("should generate different tokens for different users", func(t *testing.T) {
		token1, err1 := util.GenerateToken(1, "user1", secret)
		token2, err2 := util.GenerateToken(2, "user2", secret)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, token1, token2)
	})
}

func TestExtractTokenFromHeader(t *testing.T) {
	t.Run("should extract token from valid Bearer header", func(t *testing.T) {
		authHeader := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

		token, err := util.ExtractTokenFromHeader(authHeader)

		assert.NoError(t, err)
		assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", token)
	})

	t.Run("should handle header without Bearer prefix", func(t *testing.T) {
		authHeader := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

		token, err := util.ExtractTokenFromHeader(authHeader)

		assert.NoError(t, err)
		assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", token)
	})

	t.Run("should return error for empty header", func(t *testing.T) {
		authHeader := ""

		token, err := util.ExtractTokenFromHeader(authHeader)

		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("should return error for Bearer only header", func(t *testing.T) {
		authHeader := "Bearer "

		token, err := util.ExtractTokenFromHeader(authHeader)

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}
