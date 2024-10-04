package dto

import "time"

type CartItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gte=1"`
}

type CartItemResponse struct {
	Product    ProductDTO `json:"product" gorm:"embedded;embeddedPrefix:product_"`
	Quantity   int        `json:"quantity"`
	TotalPrice float64    `json:"total_price"`
	CreatedAt  time.Time
}
