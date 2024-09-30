package entity

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ID          string `gorm:"primaryKey"`
	Name        string
	Description string
	PriceIdr    float64
	CategoryID  string
	StoreID     string
	OrderItems  []OrderItem `gorm:"primaryKey:ProductID"`
}
