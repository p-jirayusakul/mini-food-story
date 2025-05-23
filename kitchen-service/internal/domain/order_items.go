package domain

import shareModel "food-story/shared/model"

type SearchOrderItems struct {
	Name        string
	TableNumber []int32
	StatusCode  []string
	OrderByType string
	OrderBy     string
	PageSize    int64
	PageNumber  int64
}

type SearchOrderItemsResult struct {
	TotalItems int64                    `json:"totalItems" example:"10"`
	TotalPages int64                    `json:"totalPages" example:"1"`
	Data       []*shareModel.OrderItems `json:"data"`
}
