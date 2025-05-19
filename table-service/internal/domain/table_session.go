package domain

import (
	"github.com/google/uuid"
	"time"
)

type TableSession struct {
	ID             int64
	TableID        int64
	SessionID      uuid.UUID
	NumberOfPeople int32
}

type CurrentTableSession struct {
	SessionID   uuid.UUID `json:"sessionID"`
	TableID     int64     `json:"tableID,string"`
	TableNumber int32     `json:"tableNumber"`
	Status      string    `json:"status"`
	StartedAt   time.Time `json:"startedAt"`
	OrderID     *string   `json:"orderID"`
}
