package entity

import (
	"gorm.io/gorm"
)

type Store struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Name     string
	Street   string
	City     string
	Province string
	UserID   string
	Products []Product `gorm:"foreignKey:StoreID"`
}
