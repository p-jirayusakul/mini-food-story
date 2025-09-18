package domain

import shareModel "food-story/shared/model"

type CurrentOrderItems struct {
	ID            int64   `json:"id,string" example:"1920153361642950656"`
	ProductID     int64   `json:"productID,string" example:"1920153361642950656"`
	StatusName    string  `json:"statusName" example:"กำลังเตรียมอาหาร"`
	StatusNameEN  string  `json:"statusNameEN" example:"Preparing"`
	StatusCode    string  `json:"statusCode" example:"PREPARING"`
	ProductName   string  `json:"productName" example:"ข้าวผัด"`
	ProductNameEN string  `json:"productNameEN" example:"Fried rice"`
	Price         float64 `json:"price" example:"60"`
	Quantity      int32   `json:"quantity" example:"1"`
	Note          *string `json:"note" example:"lorem ipsum"`
	CreatedAt     string  `json:"createdAt" example:"2025-05-23T11:59:50.010316+07:00"`
}

type SearchOrderItems struct {
	Name        string
	TableNumber int32
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

type SearchCurrentOrderItemsResult struct {
	TotalItems int64                `json:"totalItems" example:"10"`
	TotalPages int64                `json:"totalPages" example:"1"`
	Data       []*CurrentOrderItems `json:"data"`
}
