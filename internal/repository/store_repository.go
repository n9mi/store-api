package repository

import "store-api/internal/entity"

type StoreRepository struct {
	BaseRepository[entity.Store]
}

func NewStoreRepository() *StoreRepository {
	return new(StoreRepository)
}
