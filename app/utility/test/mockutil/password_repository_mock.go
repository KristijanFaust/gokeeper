package mockutil

import (
	"github.com/KristijanFaust/gokeeper/app/database/model"
	"github.com/stretchr/testify/mock"
	"github.com/upper/db/v4"
)

type PasswordRepositoryServiceMock struct {
	mock.Mock
}

func (service *PasswordRepositoryServiceMock) InsertNewPassword(password *model.Password) (db.InsertResult, error) {
	arguments := service.Called(password)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).(db.InsertResult), arguments.Error(1)
}

func (service *PasswordRepositoryServiceMock) UpdatePasswordById(name string, password []byte, passwordId uint64) error {
	arguments := service.Called(name, password, passwordId)
	return arguments.Error(0)
}

func (service *PasswordRepositoryServiceMock) FetchPasswordById(password *model.Password, passwordId uint64) error {
	arguments := service.Called(password, passwordId)

	if arguments.Error(0) == nil {
		password.Id = DefaultIdAsUint64
		password.UserId = DefaultIdAsUint64
		password.Name = DefaultPasswordName
		password.Password = []byte(DefaultPassword)
	}

	return arguments.Error(0)
}

func (service *PasswordRepositoryServiceMock) FetchAllByUserId(passwords *model.Passwords, userId uint64, queryFields []string) error {
	arguments := service.Called(passwords, userId, queryFields)

	if arguments.Error(0) == nil && userId == uint64(1) {
		*passwords = model.Passwords{
			model.Password{Id: uint64(1), UserId: uint64(1), Name: "Domain1", Password: []byte("Password1")},
			model.Password{Id: uint64(1), UserId: uint64(1), Name: "Domain2", Password: []byte("Password2")},
		}
	}

	return arguments.Error(0)
}

func DefaultPasswordRepositoryServiceMock() *PasswordRepositoryServiceMock {
	serviceMock := new(PasswordRepositoryServiceMock)
	serviceMock.On("InsertNewPassword", mock.Anything).Return(db.NewInsertResult(int64(1)), nil).Times(1)
	serviceMock.On("UpdatePasswordById", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
	serviceMock.On("FetchPasswordById", mock.Anything, mock.Anything).Return(nil).Times(1)
	serviceMock.On("FetchAllByUserId", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)

	return serviceMock
}
