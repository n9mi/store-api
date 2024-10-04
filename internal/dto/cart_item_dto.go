package dto

import "time"

type CartItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gte=1"`
}

type CartItemResponse struct {
	Product   ProductDTO `json:"product"`
	Quantity  int        `json:"quantity"`
	CreatedAt time.Time
}
