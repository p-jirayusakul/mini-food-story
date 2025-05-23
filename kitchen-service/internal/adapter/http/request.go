package http

type SearchOrderItems struct {
	PageNumber  int64    `query:"pageNumber" example:"1"`
	PageSize    int64    `query:"pageSize" example:"10"`
	Search      string   `query:"search" validate:"omitempty,no_special_char,max=255" example:"Rice"`
	StatusCode  []string `query:"statusCode" example:"SERVED,CANCELLED"`
	TableNumber []int32  `query:"tableNumber" example:"1,2,3"`
	OrderBy     string   `query:"orderBy" validate:"omitempty,oneof=id tableNumber statusCode productName quantity" example:"id"`
	OrderByType string   `query:"orderType" validate:"omitempty,oneof=asc desc" example:"asc"`
}

type SearchOrderItemsByOrderID struct {
	PageNumber int64 `query:"pageNumber" example:"1"`
}
