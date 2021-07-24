package repository

import (
	"github.com/KristijanFaust/gokeeper/app/database/model"
	"github.com/upper/db/v4"
)

type PasswordRepository interface {
	InsertNewPassword(password *model.Password) (db.InsertResult, error)
	FetchAllByUserId(passwords *model.Passwords, userId uint64) error
}

type passwordRepositoryService struct {
	session *db.Session
}

func NewPasswordRepositoryService(session *db.Session) *passwordRepositoryService {
	return &passwordRepositoryService{session: session}
}

func (repository *passwordRepositoryService) Password() db.Collection {
	return (*repository.session).Collection("password")
}

func (repository *passwordRepositoryService) InsertNewPassword(password *model.Password) (db.InsertResult, error) {
	return repository.Password().Insert(password)
}

func (repository *passwordRepositoryService) FetchAllByUserId(passwords *model.Passwords, userId uint64) error {
	return repository.Password().Find("user_id", userId).All(passwords)
}
