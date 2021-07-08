package model

import (
	"github.com/KristijanFaust/gokeeper/app/database"
	"github.com/upper/db/v4"
)

type User struct {
	Id       uint64 `db:"id,omitempty"`
	Email    string `db:"email"`
	Username string `db:"username"`
	Password string `db:"password"`
}

func UserCollection() db.Collection {
	return database.Session.Collection("user")
}

func (user *User) InsertNewUser() (db.InsertResult, error) {
	return UserCollection().Insert(user)
}

func (user *User) FetchByEmail(email string) error {
	return UserCollection().Find("email", email).One(user)
}
