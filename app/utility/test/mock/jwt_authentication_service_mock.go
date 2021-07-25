package mock

import "errors"

type JwtAuthenticationServiceMock struct{}

func (service *JwtAuthenticationServiceMock) GenerateJwt(userID uint64, expiredAt int64) (string, error) {
	return "", errors.New("mock generic service error")
}
