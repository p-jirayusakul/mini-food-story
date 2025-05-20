package database

import "github.com/jackc/pgx/v5/pgtype"

func (q *GetOrderWithItemsRow) GetID() int64 {
	return q.ID
}

func (q *GetOrderWithItemsRow) GetOrderID() int64 {
	return q.OrderID
}

func (q *GetOrderWithItemsRow) GetOrderNumber() string {
	return q.OrderNumber
}

func (q *GetOrderWithItemsRow) GetProductID() int64 {
	return q.ProductID
}

func (q *GetOrderWithItemsRow) GetStatusID() int64 {
	return q.StatusID
}

func (q *GetOrderWithItemsRow) GetStatusName() string {
	return q.StatusName
}

func (q *GetOrderWithItemsRow) GetStatusNameEN() string {
	return q.StatusNameEN
}

func (q *GetOrderWithItemsRow) GetStatusCode() string {
	return q.StatusCode
}

func (q *GetOrderWithItemsRow) GetProductName() string {
	return q.ProductName
}

func (q *GetOrderWithItemsRow) GetProductNameEN() string {
	return q.ProductNameEN
}

func (q *GetOrderWithItemsRow) GetPrice() pgtype.Numeric {
	return q.Price
}

func (q *GetOrderWithItemsRow) GetQuantity() int32 {
	return q.Quantity
}

func (q *GetOrderWithItemsRow) GetNote() pgtype.Text {
	return q.Note
}

func (q *GetOrderWithItemsRow) GetCreatedAt() pgtype.Timestamptz {
	return q.CreatedAt
}

func (q *GetOrderWithItemsGroupIDRow) GetID() int64 {
	return q.ID
}

func (q *GetOrderWithItemsGroupIDRow) GetOrderID() int64 {
	return q.OrderID
}

func (q *GetOrderWithItemsGroupIDRow) GetOrderNumber() string {
	return q.OrderNumber
}

func (q *GetOrderWithItemsGroupIDRow) GetProductID() int64 {
	return q.ProductID
}

func (q *GetOrderWithItemsGroupIDRow) GetStatusID() int64 {
	return q.StatusID
}

func (q *GetOrderWithItemsGroupIDRow) GetStatusName() string {
	return q.StatusName
}

func (q *GetOrderWithItemsGroupIDRow) GetStatusNameEN() string {
	return q.StatusNameEN
}

func (q *GetOrderWithItemsGroupIDRow) GetStatusCode() string {
	return q.StatusCode
}

func (q *GetOrderWithItemsGroupIDRow) GetProductName() string {
	return q.ProductName
}

func (q *GetOrderWithItemsGroupIDRow) GetProductNameEN() string {
	return q.ProductNameEN
}

func (q *GetOrderWithItemsGroupIDRow) GetPrice() pgtype.Numeric {
	return q.Price
}

func (q *GetOrderWithItemsGroupIDRow) GetQuantity() int32 {
	return q.Quantity
}

func (q *GetOrderWithItemsGroupIDRow) GetNote() pgtype.Text {
	return q.Note
}

func (q *GetOrderWithItemsGroupIDRow) GetCreatedAt() pgtype.Timestamptz {
	return q.CreatedAt
}
