package security

type PasswordSecurity interface {
	Argon2PasswordHasher
	AesPasswordCryptor
}

type PasswordSecurityService struct {
	Argon2PasswordHasher
	AesPasswordCryptor
}
