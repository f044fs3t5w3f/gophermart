package models

type User struct {
	ID           int64
	Login        string
	PasswordHash string
}

type Session struct {
	Token  string
	UserId int64
}
