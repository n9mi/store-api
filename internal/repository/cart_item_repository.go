package repository

import (
	"store-api/internal/dto"
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

func (r *CartItemRepository) FindByUserIDWithAssociation(db *gorm.DB, userID string) ([]dto.CartItemResponse, error) {
	var cartItems []dto.CartItemResponse

	if err := db.Model(new(entity.CartItem)).
		Select(`
			cart_items.quantity as quantity,
			cart_items.created_at as created_at,
			products.id as product_id,
			products.name as product_name,
			products.stock as product_stock,
			products.price_idr as product_price,
			(products.price_idr * cart_items.quantity) as total_price,
			categories.id as product_category_id,
			categories.name as product_category_name,
			stores.id as product_store_id,
			stores.name as product_store_name,
			stores.city as product_store_city`).
		Joins("inner join products on products.id = cart_items.product_id").
		Joins("inner join categories on categories.id = products.category_id").
		Joins("inner join stores on stores.id = products.store_id").
		Where("cart_items.user_id = ?", userID).
		Scan(&cartItems).Error; err != nil {
		return nil, err
	}

	return cartItems, nil
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

func (r *CartItemRepository) Empty(db *gorm.DB, userID string) error {
	return db.Where("user_id = ?", userID).
		Delete(new(entity.CartItem)).
		Error
}
