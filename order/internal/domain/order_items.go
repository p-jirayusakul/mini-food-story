package domain

type OrderItems struct {
	ID            int64   `json:"id,string"`
	OrderID       int64   `json:"orderID,string"`
	ProductID     int64   `json:"productID,string"`
	StatusID      int64   `json:"statusID,string"`
	StatusName    string  `json:"statusName"`
	StatusNameEN  string  `json:"statusNameEn"`
	ProductName   string  `json:"productName"`
	ProductNameEN string  `json:"productNameEn"`
	Price         float64 `json:"price"`
	Quantity      int32   `json:"quantity"`
	Note          *string `json:"note,omitempty"`
}

type OrderItemsStatus struct {
	ID         int64
	OrderID    int64
	StatusCode string
}
