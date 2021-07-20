package gql

import (
	"github.com/KristijanFaust/gokeeper/app/database/repository"
	"github.com/go-playground/validator"
)

type Resolver struct {
	userRepository     repository.UserRepository
	passwordRepository repository.PasswordRepository
	validator          *validator.Validate
}

func NewResolver(userRepository repository.UserRepository, passwordRepository repository.PasswordRepository) *Resolver {
	return &Resolver{userRepository: userRepository, passwordRepository: passwordRepository, validator: validator.New()}
}
