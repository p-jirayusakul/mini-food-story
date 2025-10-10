package http

type SearchMenu struct {
	PageNumber  int64   `query:"pageNumber"`
	PageSize    int64   `query:"pageSize"`
	Search      string  `query:"search" validate:"omitempty,no_special_char,max=255"`
	CategoryID  []int64 `query:"categoryID"`
	OrderBy     string  `query:"orderBy" validate:"omitempty,oneof=id name price"`
	OrderByType string  `query:"orderType" validate:"omitempty,oneof=asc desc"`
}
