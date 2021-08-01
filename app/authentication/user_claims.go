package authentication

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.StandardClaims
}
