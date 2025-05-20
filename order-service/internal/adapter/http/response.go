package http

type CurrentOrderResponse struct {
	TableNumber  int32  `json:"tableNumber"`
	StatusName   string `json:"statusName"`
	StatusNameEN string `json:"statusNameEN"`
	StatusCode   string `json:"statusCode"`
}
