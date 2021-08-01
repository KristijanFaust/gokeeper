package authentication

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
)

type userAuthentication struct {
	UserId uint64
}

var userContextKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func AuthenticationMiddleware(jwtSigningKey string) func(http.Handler) http.Handler {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			token := request.Header.Get("Authentication")
			if token == "" {
				nextHandler.ServeHTTP(writer, request)
				return
			}

			userClaims := &UserClaims{}
			decodedToken, err := DecodeJwt(token, userClaims, jwtSigningKey)
			if err != nil || !decodedToken.Valid {
				if err != nil {
					log.Printf("Error occurred while decoding JWT: %s", err)
				}
				log.Println("Invalid jwt, unauthorised request")
				writer.WriteHeader(http.StatusUnauthorized)
				nextHandler.ServeHTTP(writer, request)
				return
			}

			ctx := context.WithValue(request.Context(), userContextKey, userAuthentication{UserId: userClaims.UserID})
			request = request.WithContext(ctx)

			nextHandler.ServeHTTP(writer, request)
		})
	}
}

func DecodeJwt(token string, userClaims *UserClaims, jwtSigningKey string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, userClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSigningKey), nil
	})
}
