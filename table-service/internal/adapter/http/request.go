package http

type Table struct {
	TableNumber int32 `json:"tableNumber" validate:"required,gte=1" example:"1"`
	Seats       int32 `json:"seats" validate:"required,gte=1" example:"5"`
}

type updateTableStatus struct {
	StatusID string `json:"statusID" validate:"required" example:"1919968486671519744"`
}

type TableSession struct {
	TableID        string `json:"tableID" validate:"required" example:"1923564209627467776"`
	NumberOfPeople int32  `json:"numberOfPeople" validate:"required,gte=1" example:"3"`
}

type SearchTable struct {
	PageNumber     int64    `query:"pageNumber" example:"1"`
	PageSize       int64    `query:"pageSize" example:"10"`
	Search         string   `query:"search" validate:"omitempty,no_special_char,max=255" example:"1"`
	Seats          string   `query:"seats" validate:"omitempty,no_special_char,max=255" example:"5"`
	NumberOfPeople int32    `query:"numberOfPeople" example:"3"`
	Status         []string `query:"status" example:"1919968486671519744"`
	OrderBy        string   `query:"orderBy" validate:"omitempty,oneof=id tableNumber seats status" example:"id"`
	OrderByType    string   `query:"orderType" validate:"omitempty,oneof=asc desc" example:"asc"`
}

type SessionExtensionRequest struct {
	TableID          string `json:"tableID" validate:"required" example:"1923564209627467776"`
	RequestedMinutes int64  `json:"requestedMinutes" validate:"required,gte=1,max=120" example:"15"`
	ReasonCode       string `json:"reasonCode" validate:"required,oneof=CUSTOMER_REQUEST LATE_SERVICE SYSTEM_ERROR PARTIAL_COMP" example:"LATE_SERVICE"`
}
