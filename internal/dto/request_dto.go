package dto

type PaginateRequest struct {
	Page     int `validate:"min=1"`
	PageSize int `validate:"min=1"`
}
