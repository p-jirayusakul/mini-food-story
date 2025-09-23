package http

type OrderItems struct {
	Items []OrderItemsData `json:"items" validate:"required,gt=0,dive"`
}

type OrderItemsData struct {
	ProductID string  `json:"productID" validate:"required,gt=0" example:"1921828287366041600"`
	Quantity  int32   `json:"quantity" validate:"required,gt=0" example:"1"`
	Note      *string `json:"note" example:"lorem ipsum"`
}

type SearchOrderItemsIncomplete struct {
	PageNumber  int64    `query:"pageNumber" example:"1"`
	PageSize    int64    `query:"pageSize" example:"10"`
	Search      string   `query:"search" validate:"omitempty,no_special_char,max=255" example:""`
	StatusCode  []string `query:"statusCode" example:"PREPARING,CONFIRMED,PENDING"`
	OrderBy     string   `query:"orderBy" validate:"omitempty,oneof=id tableNumber statusCode productName quantity" example:"id"`
	OrderByType string   `query:"orderType" validate:"omitempty,oneof=asc desc" example:"asc"`
}

type SearchCurrentOrderItems struct {
	PageNumber int64 `query:"pageNumber" example:"1"`
	PageSize   int64 `query:"pageSize" example:"10"`
}
