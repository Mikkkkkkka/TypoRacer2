package model

import "time"

type User struct {
	Id              int
	Username        string
	Password        string
	Token           string
	TokenExpiration time.Time
}
