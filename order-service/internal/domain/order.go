package domain

import "github.com/google/uuid"

type Order struct {
	ID           int64  `json:"id,string"`
	TableID      int64  `json:"tableID,string"`
	TableNumber  int32  `json:"tableNumber"`
	StatusID     int64  `json:"statusID,string"`
	StatusName   string `json:"statusName"`
	StatusNameEN string `json:"statusNameEN"`
	StatusCode   string `json:"statusCode"`
}

type OrderStatus struct {
	ID         int64
	SessionID  uuid.UUID
	StatusCode string
}

type CreateOrder struct {
	SessionID  uuid.UUID
	TableID    int64
	OrderItems []OrderItems
}
