package repository

import (
	"store-api/internal/dto"
	"store-api/internal/entity"

	"gorm.io/gorm"
)

type OrderRepository struct {
	BaseRepository[entity.Order]
}

func NewOrderRepository() *OrderRepository {
	return new(OrderRepository)
}

func (r *OrderRepository) FindByOrderIDWithAssociation(db *gorm.DB, orderID string) ([]dto.OrderItem, error) {
	var orderItems []dto.OrderItem

	if err := db.Model(new(entity.OrderItem)).
		Select(`
			order_items.quantity as quantity,
			order_items.created_at as created_at,
			products.id as product_id,
			products.name as product_name,
			products.stock as product_stock,
			products.price_idr as product_price,
			(products.price_idr * order_items.quantity) as total_price,
			categories.id as product_category_id,
			categories.name as product_category_name,
			stores.id as product_store_id,
			stores.name as product_store_name,
			stores.city as product_store_city`).
		Joins("inner join products on products.id = order_items.product_id").
		Joins("inner join categories on categories.id = products.category_id").
		Joins("inner join stores on stores.id = products.store_id").
		Where("order_items.order_id = ?", orderID).
		Scan(&orderItems).Error; err != nil {
		return nil, err
	}

	return orderItems, nil
}
