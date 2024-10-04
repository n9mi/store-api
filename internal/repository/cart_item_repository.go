package repository

import (
	"store-api/internal/entity"

	"gorm.io/gorm"
)

type CartItemRepository struct {
}

func NewCartItemRepository() *CartItemRepository {
	return &CartItemRepository{}
}

func (r *CartItemRepository) FindByUserID(db *gorm.DB, userID string) ([]entity.CartItem, error) {
	var cartItems []entity.CartItem

	err := db.Where("user_id = ?", userID).Find(&cartItems).Error

	return cartItems, err
}

func (r *CartItemRepository) FindByUserIDAndProductID(db *gorm.DB, userID string, productID string) (*entity.CartItem, error) {
	var cartItem entity.CartItem

	err := db.Where("user_id = ?", userID).Where("product_id = ?", productID).First(&cartItem).Error

	return &cartItem, err
}

func (r *CartItemRepository) IsExistsByUserIDAndProductID(db *gorm.DB, userID string, productID string) (bool, error) {
	var count int

	err := db.Model(new(entity.CartItem)).
		Select("1").
		Where("user_id = ?", userID).
		Where("product_id = ?", productID).
		Scan(&count).
		Error

	return count == 1, err
}

func (r *CartItemRepository) Save(db *gorm.DB, cartItem *entity.CartItem) error {
	return db.Save(cartItem).Error
}

func (r *CartItemRepository) Delete(db *gorm.DB, userID string, productID string) error {
	return db.Where("user_id = ?", userID).
		Where("product_id = ?", productID).
		Delete(new(entity.CartItem)).
		Error
}
