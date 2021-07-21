package model

type User struct {
	Id       uint64 `db:"id,omitempty"`
	Email    string `db:"email"`
	Username string `db:"username"`
	Password []byte `db:"password"`
}
