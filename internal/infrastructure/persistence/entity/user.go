package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:string;size:20;not null;unique"`
	Password string `gorm:"type:string;size:64;not null"`
	Roles    []Role `gorm:"many2many:user_roles;"`
}
