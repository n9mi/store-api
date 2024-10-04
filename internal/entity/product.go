package entity

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ID          string `gorm:"primaryKey"`
	Name        string
	Description string
	PriceIdr    float64
	CategoryID  string
	Category    Category
	StoreID     string
	Store       Store
	Stock       int
	OrderItems  []OrderItem `gorm:"primaryKey:ProductID"`
	CartItems   []CartItem  `gorm:"primaryKey:ProductID"`
}
