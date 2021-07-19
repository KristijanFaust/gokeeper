package mock

import (
	"errors"
	"github.com/KristijanFaust/gokeeper/app/database/model"
	"github.com/upper/db/v4"
)

type PasswordRepositoryServiceMock struct{}

func (repository *PasswordRepositoryServiceMock) InsertNewPassword(password *model.Password) (db.InsertResult, error) {
	return nil, errors.New("mock generic service error")
}

func (repository *PasswordRepositoryServiceMock) FetchAllByUserId(passwords *model.Passwords, userId uint64) error {
	return errors.New("mock generic service error")
}
