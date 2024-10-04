package entity

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	ID        string `gorm:"primaryKey"`
	OrderID   string
	ProductID string
	Quantity  int
}
