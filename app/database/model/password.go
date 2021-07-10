package model

import (
	"github.com/KristijanFaust/gokeeper/app/database"
	"github.com/upper/db/v4"
)

type Password struct {
	Id       uint64 `db:"id,omitempty"`
	UserId   uint64 `db:"user_id"`
	Name     string `db:"name"`
	Password string `db:"password"`
}

type Passwords []Password

func PasswordCollection() db.Collection {
	return database.Session.Collection("password")
}

func (password *Password) InsertNewPassword() (db.InsertResult, error) {
	return PasswordCollection().Insert(password)
}

func (passwords *Passwords) FetchAllByUserId(userId uint64) error {
	return PasswordCollection().Find("user_id", userId).All(passwords)
}
