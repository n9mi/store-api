package dto

type FindAndSearchProductRequest struct {
	PaginateRequest
	QueryName          string
	QueryPriceMin      float64
	QueryPriceMax      float64
	QueryCategoryID    string
	QueryStoreID       string
	SortByPriceHighest bool
	SortByPriceLowest  bool
}

type ProductDTO struct {
	ID       string               `json:"id"`
	Name     string               `json:"name"`
	Stock    int                  `json:"stock"`
	Price    float64              `json:"price"`
	Category CategoryItemResponse `gorm:"embedded;embeddedPrefix:category_" json:"category"`
	Store    StoreItemResponse    `gorm:"embedded;embeddedPrefix:store_" json:"store"`
}
