package repository

import (
	"github.com/KristijanFaust/gokeeper/app/database/model"
	"github.com/upper/db/v4"
)

type UserRepository interface {
	InsertNewUser(user *model.User) (db.InsertResult, error)
	FetchByEmail(user *model.User, email string) error
	FetchMasterPasswordByUserId(user *model.User, id uint64) error
}

type userRepositoryService struct {
	session *db.Session
}

func NewUserRepositoryService(session *db.Session) *userRepositoryService {
	return &userRepositoryService{session: session}
}

func (userRepositoryService *userRepositoryService) User() db.Collection {
	return (*userRepositoryService.session).Collection("user")
}

func (userRepositoryService *userRepositoryService) InsertNewUser(user *model.User) (db.InsertResult, error) {
	return userRepositoryService.User().Insert(user)
}

func (userRepositoryService *userRepositoryService) FetchByEmail(user *model.User, email string) error {
	return userRepositoryService.User().Find("email", email).One(user)
}

func (userRepositoryService *userRepositoryService) FetchMasterPasswordByUserId(user *model.User, id uint64) error {
	q := (*userRepositoryService.session).SQL().Select("password").From("user").Where("id = ?", id)
	return q.One(user)
}
