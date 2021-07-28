package mockutil

import (
	"github.com/stretchr/testify/mock"
)

type JwtAuthenticationServiceMock struct {
	mock.Mock
}

func (service *JwtAuthenticationServiceMock) GenerateJwt(userID uint64, expiredAt int64) (string, error) {
	arguments := service.Called(userID, expiredAt)
	return arguments.String(0), arguments.Error(1)
}

func DefaultJwtAuthenticationServiceMock() *JwtAuthenticationServiceMock {
	serviceMock := new(JwtAuthenticationServiceMock)
	serviceMock.On("GenerateJwt", mock.Anything, mock.Anything).Return(MockedJwtToken, nil).Times(1)

	return serviceMock
}
