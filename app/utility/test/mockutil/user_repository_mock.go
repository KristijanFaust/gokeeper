package mockutil

import (
	"github.com/KristijanFaust/gokeeper/app/database/model"
	"github.com/stretchr/testify/mock"
	"github.com/upper/db/v4"
)

type UserRepositoryServiceMock struct {
	mock.Mock
}

func (service *UserRepositoryServiceMock) InsertNewUser(user *model.User) (db.InsertResult, error) {
	arguments := service.Called(user)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(db.InsertResult), arguments.Error(1)
}

func (service *UserRepositoryServiceMock) FetchByEmail(user *model.User, email string, queryFields []string) error {
	arguments := service.Called(user, email, queryFields)

	if arguments.Error(0) == nil {
		user.Id = uint64(1)
		user.Email = email
		user.Username = "username"
		user.Password = []byte(MockedUserMasterPassword)
	}

	return arguments.Error(0)
}

func (service *UserRepositoryServiceMock) FetchMasterPasswordByUserId(user *model.User, id uint64) error {
	arguments := service.Called(user, id)

	if arguments.Error(0) == nil {
		user.Password = []byte(MockedUserMasterPassword)
	}

	return arguments.Error(0)
}

func DefaultUserRepositoryServiceMock() *UserRepositoryServiceMock {
	serviceMock := new(UserRepositoryServiceMock)
	serviceMock.On("InsertNewUser", mock.Anything).Return(db.NewInsertResult(int64(1)), nil).Times(1)
	serviceMock.On("FetchByEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
	serviceMock.On("FetchMasterPasswordByUserId", mock.Anything, mock.Anything).Return(nil).Times(1)

	return serviceMock
}
