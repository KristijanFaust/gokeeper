package mock

import (
	"errors"
	"golang.org/x/crypto/argon2"
)

type PasswordSecurityServiceMock struct {
	EncryptionError bool
	DecryptionError bool
}

func (service *PasswordSecurityServiceMock) EncryptWithAes(password string, masterPassword []byte) ([]byte, error) {
	if service.EncryptionError {
		return nil, errors.New("mock generic service error")
	}

	return []byte("EncryptedPassword"), nil
}

func (service *PasswordSecurityServiceMock) DecryptWithAes(encryptedPassword []byte, masterPassword []byte) (string, error) {
	if service.DecryptionError {
		return "", errors.New("mock generic service error")
	}

	return "DecryptedPassword", nil
}

func (service *PasswordSecurityServiceMock) HashWithArgon2id(password string) []byte {
	return argon2.IDKey([]byte(password), []byte("testSalt"), 1, 1024, 1, 128)
}
