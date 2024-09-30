package entity

import "gorm.io/gorm"

type CustomerAddress struct {
	gorm.Model
	ID        string `gorm:"primaryKey"`
	Street    string
	City      string
	Province  string
	IsDefault bool
	UserID    string
	Orders    []Order `gorm:"foreignKey:CustomerAddressID"`
}
