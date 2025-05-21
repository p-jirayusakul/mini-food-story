package model

import (
	"github.com/google/uuid"
	"time"
)

type CurrentTableSession struct {
	SessionID   uuid.UUID `json:"sessionID"`
	TableID     int64     `json:"tableID,string"`
	TableNumber int32     `json:"tableNumber"`
	Status      string    `json:"status"`
	StartedAt   time.Time `json:"startedAt"`
	OrderID     *string   `json:"orderID"`
}
