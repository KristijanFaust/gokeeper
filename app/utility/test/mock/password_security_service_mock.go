package mock

import (
	"errors"
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
	return []byte("MockedHashedMasterPasswordThatIsAtLeast32BytesLong")
}
