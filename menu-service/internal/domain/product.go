package domain

type Product struct {
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

type SearchProduct struct {
	Name        string
	CategoryID  []int64
	IsAvailable bool
	OrderByType string
	OrderBy     string
	PageSize    int64
	PageNumber  int64
}

type SearchProductResult struct {
	TotalItems int64      `json:"totalItems" example:"10"`
	TotalPages int64      `json:"totalPages" example:"1"`
	Data       []*Product `json:"data"`
}
