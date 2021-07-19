package gql

import "github.com/KristijanFaust/gokeeper/app/database/repository"

type Resolver struct {
	userRepository     repository.UserRepository
	passwordRepository repository.PasswordRepository
}

func NewResolver(userRepository repository.UserRepository, passwordRepository repository.PasswordRepository) *Resolver {
	return &Resolver{userRepository: userRepository, passwordRepository: passwordRepository}
}
