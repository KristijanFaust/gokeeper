package gql

import (
	"github.com/KristijanFaust/gokeeper/app/authentication"
	"github.com/KristijanFaust/gokeeper/app/database/repository"
	"github.com/KristijanFaust/gokeeper/app/security"
	"github.com/go-playground/validator"
)

type Resolver struct {
	userRepository          repository.UserRepository
	passwordRepository      repository.PasswordRepository
	passwordSecurityService security.PasswordSecurity
	authenticationService   authentication.JwtAuthenticator
	validator               *validator.Validate
}

func NewResolver(
	userRepository repository.UserRepository,
	passwordRepository repository.PasswordRepository,
	passwordSecurityService security.PasswordSecurity,
	authenticationService authentication.JwtAuthenticator,
) *Resolver {
	return &Resolver{
		userRepository:          userRepository,
		passwordRepository:      passwordRepository,
		passwordSecurityService: passwordSecurityService,
		authenticationService:   authenticationService,
		validator:               validator.New(),
	}
}
