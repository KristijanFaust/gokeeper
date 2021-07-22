package gql

import (
	"github.com/KristijanFaust/gokeeper/app/database/repository"
	"github.com/KristijanFaust/gokeeper/app/security"
	"github.com/go-playground/validator"
)

type Resolver struct {
	userRepository        repository.UserRepository
	passwordRepository    repository.PasswordRepository
	passwordCryptoService security.AesPasswordCryptor
	validator             *validator.Validate
}

func NewResolver(
	userRepository repository.UserRepository,
	passwordRepository repository.PasswordRepository,
	passwordCryptoService security.AesPasswordCryptor,
) *Resolver {
	return &Resolver{
		userRepository:        userRepository,
		passwordRepository:    passwordRepository,
		passwordCryptoService: passwordCryptoService,
		validator:             validator.New(),
	}
}
