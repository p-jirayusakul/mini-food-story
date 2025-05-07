package http

type createTable struct {
	TableNumber int32 `json:"tableNumber" validate:"required,gte=1"`
	Seats       int32 `json:"seats" validate:"required,gte=1"`
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
