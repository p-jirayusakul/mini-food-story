package http

type OrderItems struct {
	Items []OrderItemsData `json:"items" validate:"required,gt=0,dive"`
}

type OrderItemsData struct {
	ProductID string  `json:"productID" validate:"required,gt=0"`
	Quantity  int32   `json:"quantity" validate:"required,gt=0"`
	Note      *string `json:"note"`
}

type SearchOrderItemsIncomplete struct {
	PageNumber  int64    `query:"pageNumber"`
	PageSize    int64    `query:"pageSize"`
	Search      string   `query:"search" validate:"omitempty,no_special_char,max=255"`
	StatusCode  []string `query:"statusCode"`
	OrderBy     string   `query:"orderBy" validate:"omitempty,oneof=id tableNumber statusCode productName quantity"`
	OrderByType string   `query:"orderType" validate:"omitempty,oneof=asc desc"`
}

type SearchOrderItems struct {
	PageNumber  int64  `query:"pageNumber"`
	PageSize    int64  `query:"pageSize"`
	OrderBy     string `query:"orderBy" validate:"omitempty,oneof=id tableNumber statusCode productName quantity"`
	OrderByType string `query:"orderType" validate:"omitempty,oneof=asc desc"`
}
