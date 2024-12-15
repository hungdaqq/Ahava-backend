package models

type AddToCart struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type ProductOffer struct {
	ID              uint    `json:"id"`
	CategoryID      int     `json:"category_id"`
	Image           string  `json:"image"`
	ProductName     string  `json:"product_name"`
	Size            string  `json:"size"`
	Stock           int     `json:"stock"`
	Price           float64 `json:"price"`
	DiscountedPrice float64 `json:"discounted_price"`
}

type Banner struct {
	CategoryID         int
	CategoryName       string
	DiscountPercentage int
	Images             []string `gorm:"-"`
}
