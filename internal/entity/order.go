package entity

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	ID                string `gorm:"primaryKey"`
	UserID            string
	CustomerAddressID string
	PaymentMethodID   string
	OrderItems        []OrderItem `gorm:"primaryKey:OrderID"`
}
