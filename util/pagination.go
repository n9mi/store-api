package util

import (
	"math"
	"store-api/internal/dto"

	"gorm.io/gorm"
)

func Paginate[T any](db *gorm.DB, page int, pageSize int, pagination *dto.Pagination) func(db *gorm.DB) *gorm.DB {
	_db := db

	var totalRows int64
	_db.Model(new(T)).Count(&totalRows)

	pagination.CurrentPage = page
	pagination.TotalPage = int(math.Ceil(float64(totalRows) / float64(pageSize)))
	pagination.PageSize = pageSize

	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
