package database

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
