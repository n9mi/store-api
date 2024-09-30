package repository

import "store-api/internal/entity"

type ProductRepository struct {
	BaseRepository[entity.Product]
}

func NewProductRepository() *ProductRepository {
	return new(ProductRepository)
}
