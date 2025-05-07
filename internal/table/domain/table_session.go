package domain

import "github.com/google/uuid"

type TableSession struct {
	ID             int64
	TableID        int64
	SessionID      uuid.UUID
	NumberOfPeople int32
}
