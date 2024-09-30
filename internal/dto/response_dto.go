package dto

type Response[T any] struct {
	Status     string            `json:"status"`
	Messages   map[string]string `json:"messages"`
	Data       T                 `json:"data,omitempty"`
	Pagination *Pagination       `json:"pagination,omitempty"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	PageSize    int `json:"page_size"`
}
