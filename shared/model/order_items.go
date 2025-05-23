package model

import (
	"food-story/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

type OrderItemsRow interface {
	GetID() int64
	GetOrderID() int64
	GetOrderNumber() string
	GetProductID() int64
	GetStatusID() int64
	GetStatusName() string
	GetStatusNameEN() string
	GetStatusCode() string
	GetProductName() string
	GetProductNameEN() string
	GetTableNumber() int32
	GetPrice() pgtype.Numeric
	GetQuantity() int32
	GetNote() pgtype.Text
	GetCreatedAt() pgtype.Timestamptz
}

type OrderItems struct {
	ID            int64   `json:"id,string" example:"1920153361642950656"`
	OrderID       int64   `json:"orderID,string" example:"1921828287366041600"`
	OrderNumber   string  `json:"orderNumber" example:"FS-20250523-0001"`
	ProductID     int64   `json:"productID,string" example:"1921822053405560832"`
	StatusID      int64   `json:"statusID,string" example:"1921868485739155458"`
	TableNumber   int32   `json:"tableNumber" example:"1"`
	StatusName    string  `json:"statusName" example:"กำลังเตรียมอาหาร"`
	StatusNameEN  string  `json:"statusNameEN" example:"Preparing"`
	StatusCode    string  `json:"statusCode" example:"PREPARING"`
	ProductName   string  `json:"productName" example:"ข้าวผัด"`
	ProductNameEN string  `json:"productNameEN" example:"Fried rice"`
	Price         float64 `json:"price" example:"60"`
	Quantity      int32   `json:"quantity" example:"1"`
	Note          *string `json:"note" example:"lorem ipsum"`
	CreatedAt     string  `json:"createdAt" example:"2025-05-23T13:50:36+07:00"`
}

type OrderItemsStatus struct {
	ID         int64
	OrderID    int64
	StatusCode string
}

func TransformOrderItemsResults[T OrderItemsRow](results []T) []*OrderItems {
	data := make([]*OrderItems, len(results))
	for index, row := range results {
		createdAt, _ := utils.PgTimestampToThaiISO8601(row.GetCreatedAt())
		data[index] = &OrderItems{
			ID:            row.GetID(),
			OrderID:       row.GetOrderID(),
			OrderNumber:   row.GetOrderNumber(),
			ProductID:     row.GetProductID(),
			StatusID:      row.GetStatusID(),
			TableNumber:   row.GetTableNumber(),
			StatusName:    row.GetStatusName(),
			StatusNameEN:  row.GetStatusNameEN(),
			StatusCode:    row.GetStatusCode(),
			ProductName:   row.GetProductName(),
			ProductNameEN: row.GetProductNameEN(),
			Price:         utils.PgNumericToFloat64(row.GetPrice()),
			Quantity:      row.GetQuantity(),
			Note:          utils.PgTextToStringPtr(row.GetNote()),
			CreatedAt:     createdAt,
		}
	}
	return data
}

func TransformOrderItemsByIDResults[T OrderItemsRow](results T) *OrderItems {
	createdAt, _ := utils.PgTimestampToThaiISO8601(results.GetCreatedAt())
	return &OrderItems{
		ID:            results.GetID(),
		OrderID:       results.GetOrderID(),
		OrderNumber:   results.GetOrderNumber(),
		ProductID:     results.GetProductID(),
		StatusID:      results.GetStatusID(),
		TableNumber:   results.GetTableNumber(),
		StatusName:    results.GetStatusName(),
		StatusNameEN:  results.GetStatusNameEN(),
		StatusCode:    results.GetStatusCode(),
		ProductName:   results.GetProductName(),
		ProductNameEN: results.GetProductNameEN(),
		Price:         utils.PgNumericToFloat64(results.GetPrice()),
		Quantity:      results.GetQuantity(),
		Note:          utils.PgTextToStringPtr(results.GetNote()),
		CreatedAt:     createdAt,
	}
}
