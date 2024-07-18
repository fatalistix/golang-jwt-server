package model

type User struct {
	Id       uint
	Username string
	Password string
	Roles    []Role
}
