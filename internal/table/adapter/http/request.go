package http

type Table struct {
	TableNumber int32 `json:"tableNumber" validate:"required,gte=1"`
	Seats       int32 `json:"seats" validate:"required,gte=1"`
}

type updateTableStatus struct {
	StatusID string `json:"statusID" validate:"required"`
}

type TableSession struct {
	TableID        string `json:"tableID" validate:"required"`
	NumberOfPeople int32  `json:"numberOfPeople" validate:"required,gte=1"`
}

type SearchTable struct {
	PageNumber     int64    `query:"pageNumber"`
	PageSize       int64    `query:"pageSize"`
	Search         string   `query:"search" validate:"omitempty,no_special_char,max=255"`
	Seats          string   `query:"seats" validate:"omitempty,no_special_char,max=255"`
	NumberOfPeople int32    `query:"numberOfPeople"`
	Status         []string `query:"status"`
	OrderBy        string   `query:"orderBy" validate:"omitempty,oneof=id tableNumber seats status"`
	OrderByType    string   `query:"orderType" validate:"omitempty,oneof=asc desc"`
}
