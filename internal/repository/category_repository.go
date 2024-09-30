package repository

import "store-api/internal/entity"

type CategoryRepository struct {
	BaseRepository[entity.Category]
}

func NewCategoryRepository() *CategoryRepository {
	return new(CategoryRepository)
}
