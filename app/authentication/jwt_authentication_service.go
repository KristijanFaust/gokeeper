package authentication

import (
	"context"
	"github.com/KristijanFaust/gokeeper/app/config"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

// Variable meant for mocking
var signingCall = func(token *jwt.Token, signingKey []byte) (string, error) { return token.SignedString(signingKey) }

type JwtAuthenticator interface {
	GenerateJwt(userID uint64) (string, error)
	GetAuthenticatedUserDataFromContext(context context.Context) *UserAuthentication
}

type jwtAuthenticationService struct {
	issuer        string
	jwtSigningKey []byte
	expiresAt     int64
}

func NewJwtAuthenticationService(authenticationConfig *config.Authentication) *jwtAuthenticationService {
	return &jwtAuthenticationService{
		issuer:        authenticationConfig.Issuer,
		jwtSigningKey: []byte(authenticationConfig.JwtSigningKey),
		expiresAt:     time.Now().Add(time.Minute * time.Duration(authenticationConfig.JwtDurationInMinutes)).Unix(),
	}
}

func (service *jwtAuthenticationService) GenerateJwt(userID uint64) (string, error) {
	userClaims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: service.expiresAt,
			Issuer:    service.issuer,
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	signedToken, err := signingCall(token, service.jwtSigningKey)
	if err != nil {
		log.Printf("Error occurred while generating jwt token: %s", err)
		return "", err
	}

	return signedToken, nil
}

func (service *jwtAuthenticationService) GetAuthenticatedUserDataFromContext(context context.Context) *UserAuthentication {
	return context.Value(userContextKey).(*UserAuthentication)
}
