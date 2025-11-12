package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// GetID Implement TableRow for QuickSearchTablesRow
func (q *QuickSearchTablesRow) GetID() int64 {
	return q.ID
}

func (q *QuickSearchTablesRow) GetTableNumber() int32 {
	return q.TableNumber
}

func (q *QuickSearchTablesRow) GetStatus() string {
	return q.Status
}

func (q *QuickSearchTablesRow) GetStatusEN() string {
	return q.StatusEN
}

func (q *QuickSearchTablesRow) GetStatusCode() string {
	return q.StatusCode
}

func (q *QuickSearchTablesRow) GetSeats() int32 {
	return q.Seats
}

func (q *QuickSearchTablesRow) GetOrderID() *int64 {
	if q.OrderID.Valid {
		return &q.OrderID.Int64
	}
	return nil
}

func (q *QuickSearchTablesRow) GetExpiresAt() pgtype.Timestamptz {
	return q.ExpiresAt
}

func (q *QuickSearchTablesRow) GetExtendTotalMinutes() int32 {
	if q.ExtendTotalMinutes.Valid {
		return q.ExtendTotalMinutes.Int32
	}
	return 0
}

// GetID Implement TableRow for SearchTablesRow
func (s *SearchTablesRow) GetID() int64 {
	return s.ID
}

func (s *SearchTablesRow) GetTableNumber() int32 {
	return s.TableNumber
}

func (s *SearchTablesRow) GetStatus() string {
	return s.Status
}

func (s *SearchTablesRow) GetStatusEN() string {
	return s.StatusEN
}

func (s *SearchTablesRow) GetStatusCode() string {
	return s.StatusCode
}

func (s *SearchTablesRow) GetSeats() int32 {
	return s.Seats
}

func (s *SearchTablesRow) GetOrderID() *int64 {
	if s.OrderID.Valid {
		return &s.OrderID.Int64
	}
	return nil
}

func (s *SearchTablesRow) GetExpiresAt() pgtype.Timestamptz {
	return s.ExpiresAt
}

func (s *SearchTablesRow) GetExtendTotalMinutes() int32 {
	if s.ExtendTotalMinutes.Valid {
		return s.ExtendTotalMinutes.Int32
	}
	return 0
}
