package model

type Password struct {
	Id       uint64 `db:"id,omitempty"`
	UserId   uint64 `db:"user_id"`
	Name     string `db:"name"`
	Password []byte `db:"password"`
}

type Passwords []Password
