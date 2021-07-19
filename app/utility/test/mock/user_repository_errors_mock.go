package mock

import (
	"errors"
	"github.com/KristijanFaust/gokeeper/app/database/model"
	"github.com/upper/db/v4"
)

type UserRepositoryServiceMock struct{}

func (userRepositoryService *UserRepositoryServiceMock) InsertNewUser(user *model.User) (db.InsertResult, error) {
	return nil, errors.New("mock generic service error")
}

func (userRepositoryService *UserRepositoryServiceMock) FetchByEmail(user *model.User, email string) error {
	return errors.New("mock generic service error")
}
