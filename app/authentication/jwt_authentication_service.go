package authentication

import (
	"github.com/dgrijalva/jwt-go"
	"log"
)

// Variable meant for mocking
var signingCall = func(token *jwt.Token, signingKey []byte) (string, error) { return token.SignedString(signingKey) }

type UserClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.StandardClaims
}

type JwtAuthenticator interface {
	GenerateJwt(userID uint64, expiredAt int64) (string, error)
}

type jwtAuthenticationService struct {
	issuer        string
	jwtSigningKey []byte
}

func NewJwtAuthenticationService(issuer string, jwtSigningKey []byte) *jwtAuthenticationService {
	return &jwtAuthenticationService{issuer: issuer, jwtSigningKey: jwtSigningKey}
}

func (service *jwtAuthenticationService) GenerateJwt(userID uint64, expiredAt int64) (string, error) {
	userClaims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt,
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
