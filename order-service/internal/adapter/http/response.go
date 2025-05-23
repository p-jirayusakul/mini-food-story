package http

type CurrentOrderResponse struct {
	TableNumber  int32  `json:"tableNumber" example:"1"`
	StatusName   string `json:"statusName" example:"ยืนยันออเดอร์"`
	StatusNameEN string `json:"statusNameEN" example:"Confirmed"`
	StatusCode   string `json:"statusCode" example:"CONFIRMED"`
}
