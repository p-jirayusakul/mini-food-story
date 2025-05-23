package domain

type Table struct {
	ID          int64  `json:"id,string" example:"1923564209627467776"`
	TableNumber int32  `json:"tableNumber" example:"1"`
	Status      string `json:"status" example:"สั่งอาหารแล้ว"`
	StatusEn    string `json:"statusEN" example:"Ordered"`
	Seats       int32  `json:"seats" example:"5"`
}

type TableStatus struct {
	ID       int64
	StatusID int64
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
	TotalItems int64    `json:"totalItems" example:"10"`
	TotalPages int64    `json:"totalPages" example:"1"`
	Data       []*Table `json:"data"`
}
