package domain

import shareModel "food-story/shared/model"

type CurrentOrderItems struct {
	ID            int64   `json:"id,string"`
	ProductID     int64   `json:"productID,string"`
	StatusName    string  `json:"statusName"`
	StatusNameEN  string  `json:"statusNameEN"`
	StatusCode    string  `json:"statusCode"`
	ProductName   string  `json:"productName"`
	ProductNameEN string  `json:"productNameEN"`
	Price         float64 `json:"price"`
	Quantity      int32   `json:"quantity"`
	Note          *string `json:"note"`
	CreatedAt     string  `json:"createdAt"`
}

type SearchOrderItems struct {
	Name        string
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

type SearchCurrentOrderItemsResult struct {
	TotalItems int64                `json:"totalItems"`
	TotalPages int64                `json:"totalPages"`
	Data       []*CurrentOrderItems `json:"data"`
}
