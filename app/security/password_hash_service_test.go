package security

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// HashWithArgon2id should successfully hash a given value
func TestHashWithArgon2id(t *testing.T) {
	passwordHashService := PasswordHashService{}
	hashedPassword := passwordHashService.HashWithArgon2id("TestPassword")
	fmt.Println(len(hashedPassword))
	assert.NotNil(t, hashedPassword, "Should return a hashed value")
	assert.Equal(t, len(hashedPassword), 128, "Hashed value should be of expected length")
}
