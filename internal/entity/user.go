package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID              string `gorm:"primaryKey"`
	Name            string
	Email           string `gorm:"unique"`
	Password        string
	Roles           []*Role           `gorm:"many2many:user_roles"`
	Stores          []Store           `gorm:"foreignKey:UserID"`
	CustomerAddress []CustomerAddress `gorm:"foreignKey:UserID"`
	Orders          []Order           `gorm:"foreignKey:UserID"`
}
