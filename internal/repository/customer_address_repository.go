package repository

import (
	"store-api/internal/entity"

	"gorm.io/gorm"
)

type CustomerAddressRepository struct {
	BaseRepository[entity.CustomerAddress]
}

func NewCustomerAddressRepository() *CustomerAddressRepository {
	return new(CustomerAddressRepository)
}

func (r *CustomerAddressRepository) FindByUserID(db *gorm.DB, userID string) ([]entity.CustomerAddress, error) {
	var custAddresses []entity.CustomerAddress
	err := db.Where("user_id = ?", userID).Find(&custAddresses).Error

	return custAddresses, err
}
