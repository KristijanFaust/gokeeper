package repository

import (
	"github.com/KristijanFaust/gokeeper/app/database/model"
	"github.com/iancoleman/strcase"
	"github.com/upper/db/v4"
)

type PasswordRepository interface {
	InsertNewPassword(password *model.Password) (db.InsertResult, error)
	FetchAllByUserId(passwords *model.Passwords, userId uint64, queryFields []string) error
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

func (repository *passwordRepositoryService) FetchAllByUserId(passwords *model.Passwords, userId uint64, queryFields []string) error {
	query := (*repository.session).SQL().Select().Columns()
	for _, field := range queryFields {
		query = query.Columns(strcase.ToSnake(field))
	}
	return query.From("password").Where("user_id = ?", userId).All(passwords)
}
