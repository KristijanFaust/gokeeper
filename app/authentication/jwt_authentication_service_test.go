package authentication

import (
	"context"
	"errors"
	"github.com/KristijanFaust/gokeeper/app/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

// GenerateJwt should successfully generate a json web token
func TestGenerateJwt(t *testing.T) {
	authenticationService := setupAuthenticationService()
	token, err := authenticationService.GenerateJwt(uint64(1))
	assert.Nil(t, err, "Should not return an error")
	assert.NotNil(t, token, "Jwt token should be generated")
}

// GenerateJwt should return an error in case signing fails
func TestGenerateJwtWithSigningError(t *testing.T) {
	authenticationService := setupAuthenticationService()
	signingCall = func(token *jwt.Token, signingKey []byte) (string, error) { return "", errors.New("mocked error") }
	token, err := authenticationService.GenerateJwt(uint64(1))
	assert.Equal(t, err, errors.New("mocked error"), "Should return signing error when signing fails")
	assert.Equal(t, token, "", "Jwt token should not be generated")
}

// GetAuthenticatedUserDataFromContext should successfully get user authentication details from context
func TestGetAuthenticatedUserDataFromContext(t *testing.T) {
	authenticationService := setupAuthenticationService()
	ctx := context.WithValue(context.Background(), userContextKey, &UserAuthentication{UserId: uint64(1)})
	userAuthentication := authenticationService.GetAuthenticatedUserDataFromContext(ctx)
	assert.Equal(t, userAuthentication, &UserAuthentication{UserId: uint64(1)}, "Should return signing error when signing fails")
}

// GetAuthenticatedUserDataFromContext should return nil when user authentication data is not provided
func TestGetAuthenticatedUserDataFromContextWithoutAuthenticationData(t *testing.T) {
	authenticationService := setupAuthenticationService()
	userAuthentication := authenticationService.GetAuthenticatedUserDataFromContext(context.Background())
	assert.Nil(t, userAuthentication)
}

func setupAuthenticationService() *jwtAuthenticationService {
	return NewJwtAuthenticationService(
		&config.Authentication{
			Issuer:               "issuer",
			JwtSigningKey:        "signingKey",
			JwtDurationInMinutes: 1,
		},
	)
}
