package mock

import "errors"

type PasswordCryptoServiceMock struct {
	EncryptionError bool
	DecryptionError bool
}

func (service *PasswordCryptoServiceMock) EncryptWithAes(password string, masterPassword []byte) ([]byte, error) {
	if service.EncryptionError {
		return nil, errors.New("mock generic service error")
	}

	return []byte("EncryptedPassword"), nil
}

func (service *PasswordCryptoServiceMock) DecryptWithAes(encryptedPassword []byte, masterPassword []byte) (string, error) {
	if service.DecryptionError {
		return "", errors.New("mock generic service error")
	}

	return "DecryptedPassword", nil
}
