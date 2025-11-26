package domain

import (
	"time"

	"github.com/google/uuid"
)

type ListSessionExtensionReason struct {
	ID       int64   `json:"id,string"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	NameEN   string  `json:"nameEN"`
	Category *string `json:"category"`
	ModeCode *string `json:"modeCode"`
}

type TableSession struct {
	ID             int64
	TableID        int64
	SessionID      uuid.UUID
	NumberOfPeople int32
}

type SessionExtension struct {
	TableID    int64
	ProductID  int64
	ReasonCode string
}

type CurrentTableSession struct {
	SessionID   uuid.UUID `json:"sessionID" example:"a9213539-b135-42cc-b714-60cfd1b099ec"`
	TableID     int64     `json:"tableID,string" example:"1920153361642950656"`
	TableNumber int32     `json:"tableNumber" example:"1"`
	Status      string    `json:"status" example:"active"`
	StartedAt   time.Time `json:"startedAt" example:"2025-05-23T11:59:50.010316+07:00"`
	OrderID     *string   `json:"orderID" example:"1922535048335069184"`
}

type ProductTimeExtension struct {
	ID             int64   `json:"id,string" example:"1921144250070732800"`
	Name           string  `json:"name" example:"ข้าวมันไก่"`
	NameEN         string  `json:"nameEN" example:"Chicken rice"`
	CategoryName   string  `json:"categoryName" example:"อาหาร"`
	CategoryNameEN string  `json:"categoryNameEN" example:"Food"`
	CategoryID     int64   `json:"categoryID,string" example:"1921143886227443712"`
	Price          float64 `json:"price" example:"100"`
	Description    *string `json:"description" example:"lorem ipsum"`
	IsAvailable    bool    `json:"isAvailable" example:"true"`
	ImageURL       *string `json:"imageURL" example:"https://example.com/image.jpg"`
}
