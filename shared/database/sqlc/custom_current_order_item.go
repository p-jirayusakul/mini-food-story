package database

import "github.com/jackc/pgx/v5/pgtype"

// ---- GetOrderWithItemsRow ---- /

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

func (q *GetOrderWithItemsRow) GetImageURL() pgtype.Text {
	return q.ImageURL
}

func (q *GetOrderWithItemsRow) GetTableNumber() int32 {
	return q.TableNumber
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
