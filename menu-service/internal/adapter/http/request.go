package http

type SearchMenu struct {
	PageNumber  int64   `query:"pageNumber"`
	PageSize    int64   `query:"pageSize"`
	Search      string  `query:"search" validate:"omitempty,no_special_char,max=255"`
	CategoryID  []int64 `query:"categoryID"`
	IsAvailable bool    `query:"isAvailable"`
	OrderBy     string  `query:"orderBy" validate:"omitempty,oneof=id tableNumber seats status"`
	OrderByType string  `query:"orderType" validate:"omitempty,oneof=asc desc"`
}
