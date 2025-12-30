package unit

import (
	"testing"

	"go_api/internal/util"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	t.Run("should hash password successfully", func(t *testing.T) {
		password := "testpassword123"
		hashedPassword, err := util.HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
		assert.NotEqual(t, password, hashedPassword)
	})

	t.Run("should produce different hashes for same password", func(t *testing.T) {
		password := "testpassword123"
		hash1, err1 := util.HashPassword(password)
		hash2, err2 := util.HashPassword(password)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2) // bcrypt produces different hashes each time
	})
}

func TestComparePassword(t *testing.T) {
	t.Run("should return true for matching password", func(t *testing.T) {
		password := "testpassword123"
		hashedPassword, _ := util.HashPassword(password)

		result := util.ComparePassword(password, hashedPassword)

		assert.True(t, result)
	})

	t.Run("should return false for non-matching password", func(t *testing.T) {
		password := "testpassword123"
		wrongPassword := "wrongpassword"
		hashedPassword, _ := util.HashPassword(password)

		result := util.ComparePassword(wrongPassword, hashedPassword)

		assert.False(t, result)
	})

	t.Run("should return false for empty password", func(t *testing.T) {
		hashedPassword, _ := util.HashPassword("somepassword")

		result := util.ComparePassword("", hashedPassword)

		assert.False(t, result)
	})
}
