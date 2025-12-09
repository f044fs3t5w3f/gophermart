package models

import "time"

type User struct {
	Id           int64
	Login        string
	PasswordHash string
}

type Session struct {
	Token  string
	UserId int64
}

const (
	StatusNew        = "NEW"
	StatusProcessing = "PROCESSING"
	StatusInvalid    = "INVALID"
	StatusProcessed  = "PROCESSED"
)

type Order struct {
	Id         int64
	UserId     int64
	Status     string
	Number     string
	UploadedAt time.Time
}
