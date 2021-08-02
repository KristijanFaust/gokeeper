package mockutil

import (
	"context"
	"github.com/KristijanFaust/gokeeper/app/authentication"
	"github.com/stretchr/testify/mock"
)

type JwtAuthenticationServiceMock struct {
	mock.Mock
}

func (service *JwtAuthenticationServiceMock) GenerateJwt(userID uint64) (string, error) {
	arguments := service.Called(userID)
	return arguments.String(0), arguments.Error(1)
}

func (service *JwtAuthenticationServiceMock) GetAuthenticatedUserDataFromContext(context context.Context) *authentication.UserAuthentication {
	arguments := service.Called(context)
	if arguments.Get(0) == nil {
		return nil
	}

	return arguments.Get(0).(*authentication.UserAuthentication)
}

func DefaultJwtAuthenticationServiceMock() *JwtAuthenticationServiceMock {
	serviceMock := new(JwtAuthenticationServiceMock)
	serviceMock.On("GenerateJwt", mock.Anything, mock.Anything).Return(MockedJwtToken, nil).Times(1)
	serviceMock.On("GetAuthenticatedUserDataFromContext", mock.Anything).Return(
		&authentication.UserAuthentication{UserId: DefaultIdAsUint64},
	).Times(1)

	return serviceMock
}
