package domain

import (
	"github.com/google/uuid"
	"time"
)

type TableSession struct {
	ID             string
	TableID        string
	SessionID      uuid.UUID
	NumberOfPeople int32
}

type CurrentTableSession struct {
	SessionID   uuid.UUID `json:"sessionID"`
	TableID     string    `json:"tableID"`
	TableNumber int32     `json:"tableNumber"`
	Status      string    `json:"status"`
	StartedAt   time.Time `json:"startedAt"`
	OrderID     *string   `json:"orderID"`
}
