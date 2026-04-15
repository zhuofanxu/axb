package dto

type PaginationResponse[T any] struct {
	Total   int64 `json:"total"`
	Page    int   `json:"pageIndex"`
	Size    int   `json:"pageSize"`
	Records []T   `json:"data"`
}

type PaginationParam struct {
	Page     int `json:"pageIndex" binding:"omitempty,min=0"`
	PageSize int `json:"pageSize" binding:"omitempty,min=1,max=1000"`
}
