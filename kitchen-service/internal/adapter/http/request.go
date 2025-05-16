package http

type SearchOrderItems struct {
	PageNumber  int64    `query:"pageNumber"`
	PageSize    int64    `query:"pageSize"`
	Search      string   `query:"search" validate:"omitempty,no_special_char,max=255"`
	StatusCode  []string `query:"statusCode"`
	TableNumber []int32  `query:"tableNumber"`
	OrderBy     string   `query:"orderBy" validate:"omitempty,oneof=id tableNumber statusCode productName quantity"`
	OrderByType string   `query:"orderType" validate:"omitempty,oneof=asc desc"`
}
