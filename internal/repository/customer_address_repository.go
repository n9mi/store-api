package repository

import "store-api/internal/entity"

type CustomerAddressRepository struct {
	BaseRepository[entity.CustomerAddress]
}

func NewCustomerAddressRepository() *CustomerAddressRepository {
	return new(CustomerAddressRepository)
}
