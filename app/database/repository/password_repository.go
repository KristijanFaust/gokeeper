package repository

import (
	"github.com/KristijanFaust/gokeeper/app/database"
	"github.com/KristijanFaust/gokeeper/app/database/model"
	"github.com/upper/db/v4"
)

type PasswordRepository interface {
	InsertNewPassword(password *model.Password) (db.InsertResult, error)
	FetchAllByUserId(passwords *model.Passwords, userId uint64) error
}

type PasswordRepositoryService struct{}

func PasswordCollection() db.Collection {
	return database.Session.Collection("password")
}

func (repository *PasswordRepositoryService) InsertNewPassword(password *model.Password) (db.InsertResult, error) {
	return PasswordCollection().Insert(password)
}

func (repository *PasswordRepositoryService) FetchAllByUserId(passwords *model.Passwords, userId uint64) error {
	return PasswordCollection().Find("user_id", userId).All(passwords)
}
