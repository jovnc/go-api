package unit

import (
	"testing"

	"go_api/internal/util"

	"github.com/stretchr/testify/assert"
)

// Test struct with validation tags
type TestUser struct {
	Email    string `validate:"required,email"`
	Username string `validate:"required,min=3,max=20"`
	Age      int    `validate:"required,gte=0,lte=130"`
}

func TestValidate(t *testing.T) {
	t.Run("should pass validation for valid struct", func(t *testing.T) {
		user := TestUser{
			Email:    "test@example.com",
			Username: "testuser",
			Age:      25,
		}

		err := util.Validate(user)

		assert.NoError(t, err)
	})

	t.Run("should fail validation for missing required field", func(t *testing.T) {
		user := TestUser{
			Email: "test@example.com",
			Age:   25,
			// Username is missing
		}

		err := util.Validate(user)

		assert.Error(t, err)
	})

	t.Run("should fail validation for invalid email", func(t *testing.T) {
		user := TestUser{
			Email:    "invalid-email",
			Username: "testuser",
			Age:      25,
		}

		err := util.Validate(user)

		assert.Error(t, err)
	})

	t.Run("should fail validation for username too short", func(t *testing.T) {
		user := TestUser{
			Email:    "test@example.com",
			Username: "ab", // min is 3
			Age:      25,
		}

		err := util.Validate(user)

		assert.Error(t, err)
	})

	t.Run("should fail validation for age out of range", func(t *testing.T) {
		user := TestUser{
			Email:    "test@example.com",
			Username: "testuser",
			Age:      150, // max is 130
		}

		err := util.Validate(user)

		assert.Error(t, err)
	})
}
