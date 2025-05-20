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
	TotalItems int64                    `json:"totalItems"`
	TotalPages int64                    `json:"totalPages"`
	Data       []*shareModel.OrderItems `json:"data"`
}
