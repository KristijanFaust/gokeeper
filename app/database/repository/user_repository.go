package repository

import (
	"github.com/KristijanFaust/gokeeper/app/database/model"
	"github.com/iancoleman/strcase"
	"github.com/upper/db/v4"
)

type UserRepository interface {
	InsertNewUser(user *model.User) (db.InsertResult, error)
	FetchByEmail(user *model.User, email string, queryFields []string) error
	FetchMasterPasswordByUserId(user *model.User, id uint64) error
}

type userRepositoryService struct {
	session *db.Session
}

func NewUserRepositoryService(session *db.Session) *userRepositoryService {
	return &userRepositoryService{session: session}
}

func (repository *userRepositoryService) User() db.Collection {
	return (*repository.session).Collection("user")
}

func (repository *userRepositoryService) InsertNewUser(user *model.User) (db.InsertResult, error) {
	return repository.User().Insert(user)
}

func (repository *userRepositoryService) FetchByEmail(user *model.User, email string, queryFields []string) error {
	query := (*repository.session).SQL().Select().Columns()
	for _, field := range queryFields {
		query = query.Columns(strcase.ToSnake(field))
	}
	return query.From("user").Where("email = ?", email).One(user)
}

func (repository *userRepositoryService) FetchMasterPasswordByUserId(user *model.User, id uint64) error {
	return (*repository.session).SQL().Select("password").From("user").Where("id = ?", id).One(user)
}
