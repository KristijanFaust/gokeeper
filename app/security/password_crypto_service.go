package security

import (
	"crypto/aes"
	"crypto/cipher"
)

// Beware that changing these constants will break compatibility with old hashed values
const (
	encryptSalt = "r95Ai4Ubur6ZXE6C"
	keyByteSize = 32
)

// Variables meant for mocking
var (
	generateNewCipherBlock = aes.NewCipher
	wrapBlockWithGCM       = cipher.NewGCM
)

type AesPasswordCryptor interface {
	EncryptWithAes(password string, masterPassword []byte) ([]byte, error)
	DecryptWithAes(encryptedPassword []byte, masterPassword []byte) (string, error)
}

type PasswordCryptoService struct{}

func (service *PasswordCryptoService) EncryptWithAes(password string, masterPassword []byte) ([]byte, error) {
	gcm, err := setUpAes(masterPassword[:keyByteSize])
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	encryptedPassword := gcm.Seal(nonce, nonce, []byte(password), []byte(encryptSalt))

	return encryptedPassword, nil
}

func (service *PasswordCryptoService) DecryptWithAes(encryptedPassword []byte, masterPassword []byte) (string, error) {
	gcm, err := setUpAes(masterPassword[:keyByteSize])
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	decryptedPassword, err := gcm.Open(nil, encryptedPassword[:nonceSize], encryptedPassword[nonceSize:], []byte(encryptSalt))
	if err != nil {
		return "", err
	}

	return string(decryptedPassword), nil
}

func setUpAes(key []byte) (cipher.AEAD, error) {
	block, err := generateNewCipherBlock(key)
	if err != nil {
		return nil, err
	}

	gcm, err := wrapBlockWithGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm, nil
}
