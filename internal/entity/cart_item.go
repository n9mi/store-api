package entity

import "time"

type CartItem struct {
	UserID    string `gorm:"primaryKey"`
	ProductID string `gorm:"primaryKey"`
	Quantity  int
	CreatedAt time.Time
}
