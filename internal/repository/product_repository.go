package repository

import (
	"store-api/internal/dto"
	"store-api/internal/entity"
	"store-api/util"
	"strings"

	"gorm.io/gorm"
)

type ProductRepository struct {
	BaseRepository[entity.Product]
}

func NewProductRepository() *ProductRepository {
	return new(ProductRepository)
}

func (r *ProductRepository) FindAll(db *gorm.DB, request dto.FindAndSearchProductRequest) ([]dto.ProductDTO, *dto.Pagination, error) {
	var productItems []dto.ProductDTO

	db = db.Model(new(entity.Product)).
		Select(`products.id as id,
			products.name as name,
			products.stock as stock,
			products.price_idr as price,
			categories.id as category_id,
			categories.name as category_name,
			stores.id as store_id,
			stores.name as store_name,
			stores.city as store_city`).
		Joins("inner join categories on categories.id = products.category_id").
		Joins("inner join stores on stores.id = products.store_id")

	if len(request.QueryName) > 0 {
		db = db.Where("lower(products.name) like ?", "%"+strings.ToLower(request.QueryName)+"%")
	}

	if request.QueryPriceMin > 0 {
		db = db.Where("products.price >= ?", request.QueryPriceMin)
	}
	if request.QueryPriceMax > 0 {
		db = db.Where("products.price <= ?", request.QueryPriceMax)
	}

	if len(request.QueryCategoryID) > 0 {
		db = db.Where("products.category_id = ?", request.QueryCategoryID)
	}

	if len(request.QueryStoreID) > 0 {
		db = db.Where("products.store_id = ?", request.QueryStoreID)
	}

	if request.SortByPriceHighest {
		db = db.Order("products.price_idr desc")
	} else if request.SortByPriceLowest {
		db = db.Order("products.price_idr asc")
	} else {
		db = db.Order("products.created_at desc")
	}

	if request.Page < 1 && request.PageSize < 1 {
		// default
		request.Page = 1
		request.PageSize = 10
	}
	var paginationRes dto.Pagination
	db = db.Scopes(util.Paginate[entity.Product](db, request.Page, request.PageSize, &paginationRes))

	if err := db.Scan(&productItems).Error; err != nil {
		return nil, nil, err
	}

	return productItems, &paginationRes, nil
}
