package domain

type Payment struct {
	ID      int64
	OrderID int64
	Method  int64
	Note    *string
}
