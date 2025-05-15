package http

type OrderItems struct {
	Items []OrderItemsData `json:"items" validate:"required,gt=0,dive"`
}

type OrderItemsData struct {
	ProductID string  `json:"productID" validate:"required,gt=0"`
	Quantity  int32   `json:"quantity" validate:"required,gt=0"`
	Note      *string `json:"note"`
}
