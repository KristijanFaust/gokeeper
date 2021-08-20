package model

type NewUser struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=1,max=32"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type NewPassword struct {
	UserID   string `json:"userId" validate:"required"`
	Name     string `json:"name" validate:"required,min=1,max=64"`
	Password string `json:"password" validate:"required"`
}

type UpdatePassword struct {
	ID       string `json:"id" validate:"required"`
	UserID   string `json:"userId" validate:"required"`
	Name     string `json:"name" validate:"required,min=1,max=64"`
	Password string `json:"password" validate:"required"`
}
