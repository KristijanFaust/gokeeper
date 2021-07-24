package mock

import (
	"errors"
	"github.com/KristijanFaust/gokeeper/app/database/model"
	"github.com/upper/db/v4"
)

type UserRepositoryServiceMock struct {
	InsertNewUserError               bool
	FetchByEmailError                bool
	FetchMasterPasswordByUserIdError bool
}

func (userRepositoryService *UserRepositoryServiceMock) InsertNewUser(user *model.User) (db.InsertResult, error) {
	if userRepositoryService.InsertNewUserError {
		return nil, errors.New("mock generic service error")
	}

	return db.NewInsertResult(int64(1)), nil
}

func (userRepositoryService *UserRepositoryServiceMock) FetchByEmail(user *model.User, email string) error {
	if userRepositoryService.FetchByEmailError {
		return errors.New("mock generic service error")
	}

	return nil
}

func (userRepositoryService *UserRepositoryServiceMock) FetchMasterPasswordByUserId(user *model.User, id uint64) error {
	if userRepositoryService.FetchMasterPasswordByUserIdError {
		return errors.New("mock generic service error")
	}

	user.Password = []byte("MockedHashedMasterPasswordThatIsAtLeast32BytesLong")
	return nil
}
