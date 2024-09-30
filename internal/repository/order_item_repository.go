package repository

import "store-api/internal/entity"

type OrderItemRepository struct {
	BaseRepository[entity.OrderItem]
}

func NewOrderItemRepository() *OrderItemRepository {
	return new(OrderItemRepository)
}
