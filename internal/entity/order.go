package entity

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID                   string `gorm:"primaryKey"`
	UserID               string
	CustomerAddressID    string
	PaymentMethodID      string
	PaymentCode          string
	PaymentCodeExpiredAt time.Time
	OrderItems           []OrderItem `gorm:"primaryKey:OrderID"`
}
