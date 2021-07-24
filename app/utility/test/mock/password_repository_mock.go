package mock

import (
	"errors"
	"github.com/KristijanFaust/gokeeper/app/database/model"
	"github.com/upper/db/v4"
)

type PasswordRepositoryServiceMock struct {
	InsertPasswordError   bool
	FetchAllByUserIdError bool
}

func (repository *PasswordRepositoryServiceMock) InsertNewPassword(password *model.Password) (db.InsertResult, error) {
	if repository.InsertPasswordError {
		return nil, errors.New("mock generic service error")
	}

	return db.NewInsertResult(int64(1)), nil
}

func (repository *PasswordRepositoryServiceMock) FetchAllByUserId(passwords *model.Passwords, userId uint64, queryFields []string) error {
	if repository.FetchAllByUserIdError {
		return errors.New("mock generic service error")
	}

	*passwords = model.Passwords{
		model.Password{Id: uint64(1), UserId: uint64(1), Name: "Domain1", Password: []byte("Password1")},
		model.Password{Id: uint64(1), UserId: uint64(1), Name: "Domain2", Password: []byte("Password2")},
	}

	return nil
}
