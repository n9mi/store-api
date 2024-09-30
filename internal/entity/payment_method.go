package entity

import "gorm.io/gorm"

type PaymentMethod struct {
	gorm.Model
	ID     string `gorm:"primaryKey"`
	Code   string
	Name   string
	Orders []Order `gorm:"foreignKey:PaymentMethodID"`
}
