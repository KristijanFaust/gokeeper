package security

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

const validEncryptionKey = "keyThatNeedsToBeAtLeast32BytesLong"
const mockedErrorMessage = "error"

// EncryptWithAes should successfully encrypt a given value
func TestEncryptWithAes(t *testing.T) {
	passwordCryptoService := PasswordCryptoService{}
	encryptedPassword, err := passwordCryptoService.EncryptWithAes("TestPassword", []byte(validEncryptionKey))
	assert.Nil(t, err, "Should not return any errors")
	assert.NotNil(t, encryptedPassword, "Should return an encrypted value")
}

// DecryptWithAes should successfully decrypt a given value
func TestDecryptWithAes(t *testing.T) {
	passwordCryptoService := PasswordCryptoService{}
	encryptedPassword, err := passwordCryptoService.EncryptWithAes("TestPassword", []byte(validEncryptionKey))
	decryptedPassword, err := passwordCryptoService.DecryptWithAes(encryptedPassword, []byte(validEncryptionKey))
	assert.Nil(t, err, "Should not return any errors")
	assert.Equal(t, decryptedPassword, "TestPassword")
}

// In practice with the current setup errors should never happen while encrypting or decrypting values, but here we test those scenarios just in case

// EncryptWithAes should return error when generating new cipher block fails
func TestEncryptWithAesWithCipherBlockCreationError(t *testing.T) {
	generateNewCipherBlock = func(key []byte) (cipher.Block, error) { return nil, errors.New(mockedErrorMessage) }
	defer func() { generateNewCipherBlock = aes.NewCipher }()

	passwordCryptoService := PasswordCryptoService{}
	encryptedPassword, err := passwordCryptoService.EncryptWithAes("test", []byte(validEncryptionKey))
	assert.Equal(t, err, errors.New(mockedErrorMessage), "Should return root error")
	assert.Nil(t, encryptedPassword, "Should not return any encrypted value")
}

// EncryptWithAes should return error when wrapping the generated cipher block with GCM fails
func TestEncryptWithAesWithCipherBlockWrappingGCMError(t *testing.T) {
	wrapBlockWithGCM = func(cipher cipher.Block) (cipher.AEAD, error) { return nil, errors.New(mockedErrorMessage) }
	defer func() { wrapBlockWithGCM = cipher.NewGCM }()

	passwordCryptoService := PasswordCryptoService{}
	encryptedPassword, err := passwordCryptoService.EncryptWithAes("test", []byte(validEncryptionKey))
	assert.Equal(t, err, errors.New(mockedErrorMessage), "Should return root error")
	assert.Nil(t, encryptedPassword, "Should not return any encrypted value")
}

// DecryptWithAes should return error when generating new cipher block fails
func TestDecryptWithAesWithCipherBlockCreationError(t *testing.T) {
	generateNewCipherBlock = func(key []byte) (cipher.Block, error) { return nil, errors.New(mockedErrorMessage) }
	defer func() { generateNewCipherBlock = aes.NewCipher }()

	passwordCryptoService := PasswordCryptoService{}
	decryptedPassword, err := passwordCryptoService.DecryptWithAes([]byte("EncryptedPassword"), []byte(validEncryptionKey))
	assert.Equal(t, err, errors.New(mockedErrorMessage), "Should return root error")
	assert.Equal(t, decryptedPassword, "", "Should return empty string")
}

// DecryptWithAes should return error when wrapping the generated cipher block with GCM fails
func TestDecryptWithAesWithCipherBlockWrappingGCMError(t *testing.T) {
	wrapBlockWithGCM = func(cipher cipher.Block) (cipher.AEAD, error) { return nil, errors.New(mockedErrorMessage) }
	defer func() { wrapBlockWithGCM = cipher.NewGCM }()

	passwordCryptoService := PasswordCryptoService{}
	decryptedPassword, err := passwordCryptoService.DecryptWithAes([]byte("EncryptedPassword"), []byte(validEncryptionKey))
	assert.Equal(t, err, errors.New(mockedErrorMessage), "Should return root error")
	assert.Equal(t, decryptedPassword, "", "Should return empty string")
}

// DecryptWithAes should return error on decryption failed
func TestDecryptWithAesWithDecryptionError(t *testing.T) {
	passwordCryptoService := PasswordCryptoService{}
	decryptedPassword, err := passwordCryptoService.DecryptWithAes([]byte("non-encrypted-password"), []byte(validEncryptionKey))
	assert.Equal(t, err, errors.New("cipher: message authentication failed"), "Should return root error")
	assert.Equal(t, decryptedPassword, "", "Should return empty string")
}
