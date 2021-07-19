package repository

import (
	"github.com/KristijanFaust/gokeeper/app/database"
	"github.com/KristijanFaust/gokeeper/app/database/model"
	"github.com/upper/db/v4"
)

type UserRepository interface {
	InsertNewUser(user *model.User) (db.InsertResult, error)
	FetchByEmail(user *model.User, email string) error
}

type UserRepositoryService struct{}

func UserCollection() db.Collection {
	return database.Session.Collection("user")
}

func (userRepositoryService *UserRepositoryService) InsertNewUser(user *model.User) (db.InsertResult, error) {
	return UserCollection().Insert(user)
}

func (userRepositoryService *UserRepositoryService) FetchByEmail(user *model.User, email string) error {
	return UserCollection().Find("email", email).One(user)
}
