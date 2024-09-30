package repository

import "store-api/internal/entity"

type OrderRepository struct {
	BaseRepository[entity.Order]
}

func NewOrderRepository() *OrderRepository {
	return new(OrderRepository)
}
