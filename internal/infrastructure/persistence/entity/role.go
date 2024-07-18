package entity

type Role struct {
	Name  string `gorm:"type:string;size:10;not null;unique;primaryKey;"`
	Users []User `gorm:"many2many:user_roles;"`
}
