package dto

import "time"

type OrderCreateRequest struct {
	CustomerAddressID string `json:"customer_address_id" validate:"required,min=1"`
	PaymentMethodID   string `json:"payment_method_id" validate:"required,min=1"`
}

type OrderCreateResponse struct {
	PaymentCode          string      `json:"payment_code"`
	PaymentCodeExpiredAt time.Time   `json:"payment_code_expired_at"`
	OrderItems           []OrderItem `json:"order_items"`
	TotalPrice           float64     `json:"total_price"`
}

type OrderItem struct {
	Product    ProductDTO `json:"product" gorm:"embedded;embeddedPrefix:product_"`
	Quantity   int        `json:"quantity"`
	TotalPrice float64    `json:"total_price"`
	CreatedAt  time.Time
}
