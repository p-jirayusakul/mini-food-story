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
	ID            int64   `json:"id,string"`
	OrderID       int64   `json:"orderID,string"`
	OrderNumber   string  `json:"orderNumber"`
	ProductID     int64   `json:"productID,string"`
	StatusID      int64   `json:"statusID,string"`
	TableNumber   int32   `json:"tableNumber"`
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
