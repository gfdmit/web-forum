package model

type User struct {
	ID           int    `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"-"`
}

type CreateUserInput struct {
	Login        string
	PasswordHash string
}
