// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Password struct {
	ID       string `json:"id"`
	UserID   string `json:"userId"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}