package mockutil

import (
	"github.com/stretchr/testify/mock"
)

type PasswordSecurityServiceMock struct {
	mock.Mock
}

func (service *PasswordSecurityServiceMock) EncryptWithAes(password string, masterPassword []byte) ([]byte, error) {
	arguments := service.Called(password, masterPassword)

	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}

	return arguments.Get(0).([]byte), arguments.Error(1)
}

func (service *PasswordSecurityServiceMock) DecryptWithAes(encryptedPassword []byte, masterPassword []byte) (string, error) {
	arguments := service.Called(encryptedPassword, masterPassword)
	return arguments.String(0), arguments.Error(1)
}

func (service *PasswordSecurityServiceMock) HashWithArgon2id(password string) []byte {
	arguments := service.Called(password)
	return arguments.Get(0).([]byte)
}

func DefaultPasswordSecurityServiceMock() *PasswordSecurityServiceMock {
	serviceMock := new(PasswordSecurityServiceMock)
	serviceMock.On("EncryptWithAes", mock.Anything, mock.Anything).Return([]byte(MockedEncryptedPassword), nil).Times(1)
	serviceMock.On("DecryptWithAes", mock.Anything, mock.Anything).Return(MockedDecryptedPassword, nil).Times(1)
	serviceMock.On("HashWithArgon2id", mock.Anything).Return([]byte(MockedUserMasterPassword)).Times(1)

	return serviceMock
}
