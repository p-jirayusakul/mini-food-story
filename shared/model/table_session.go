package model

import (
	"time"

	"github.com/google/uuid"
)

type CurrentTableSession struct {
	SessionID   uuid.UUID `json:"sessionID" example:"a9213539-b135-42cc-b714-60cfd1b099ec"`
	TableID     int64     `json:"tableID,string" example:"1920153361642950656"`
	TableNumber int32     `json:"tableNumber" example:"1"`
	Status      string    `json:"status" example:"active"`
	StartedAt   time.Time `json:"startedAt" example:"2025-05-23T11:59:50.010316+07:00"`
	OrderID     *string   `json:"orderID" example:"1922535048335069184"`
}
