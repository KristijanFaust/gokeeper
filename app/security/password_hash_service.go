package security

import "golang.org/x/crypto/argon2"

// Beware that changing these constants will break compatibility with old hashed values
const (
	hashSalt   = "57GUAhLmUPeJuW88"
	iterations = 8
	memory     = 8 * 1024
	threads    = 1
	keyLength  = 128
)

type Argon2PasswordHasher interface {
	HashWithArgon2id(password string) []byte
}

type PasswordHashService struct{}

func (service *PasswordHashService) HashWithArgon2id(password string) []byte {
	return argon2.IDKey([]byte(password), []byte(hashSalt), iterations, memory, threads, keyLength)
}
