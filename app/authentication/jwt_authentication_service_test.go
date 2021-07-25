package authentication

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

// GenerateJwt should successfully generate a json web token
func TestGenerateJwt(t *testing.T) {
	authenticationService := NewJwtAuthenticationService("issuer", []byte("signingKey"))
	token, err := authenticationService.GenerateJwt(uint64(1), int64(1))
	assert.Nil(t, err, "Should not return an error")
	assert.NotNil(t, token, "Jwt token should be generated")
}

// GenerateJwt should return an error in case signing fails
func TestGenerateJwtWithSigningError(t *testing.T) {
	authenticationService := NewJwtAuthenticationService("issuer", []byte("signingKey"))
	signingCall = func(token *jwt.Token, signingKey []byte) (string, error) { return "", errors.New("mocked error") }
	token, err := authenticationService.GenerateJwt(uint64(1), int64(1))
	assert.Equal(t, err, errors.New("mocked error"), "Should return signing error when signing fails")
	assert.Equal(t, token, "", "Jwt token should not be generated")
}
