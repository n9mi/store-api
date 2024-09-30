package entity

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	ID    string `gorm:"primaryKey"`
	Name  string
	Users []*User `gorm:"many2many:user_roles"`
}
