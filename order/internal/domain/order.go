package domain

import "github.com/google/uuid"

type Order struct {
	ID           int64     `json:"id,string"`
	SessionID    uuid.UUID `json:"sessionID"`
	TableID      int64     `json:"tableID,string"`
	StatusID     int64     `json:"statusID,string"`
	StatusName   string    `json:"statusName"`
	StatusNameEN string    `json:"statusNameEN"`
}

type OrderStatus struct {
	ID         int64
	SessionID  uuid.UUID
	StatusCode string
}
