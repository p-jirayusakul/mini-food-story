package domain

type Table struct {
	ID          string `json:"id"`
	TableNumber int32  `json:"tableNumber"`
	Status      string `json:"status"`
	StatusEn    string `json:"statusEN"`
	Seats       int32  `json:"seats"`
}

type TableStatus struct {
	ID       string
	StatusID string
}

type SearchTables struct {
	NumberOfPeople int32
	TableNumber    *int32
	Seats          *int32
	StatusCode     []string
	OrderByType    string
	OrderBy        string
	PageSize       int64
	PageNumber     int64
}

type SearchTablesResult struct {
	TotalItems int64    `json:"totalItems"`
	TotalPages int64    `json:"totalPages"`
	Data       []*Table `json:"data"`
}
