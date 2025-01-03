package models

type ProductOffer struct {
	ID              uint   `json:"id"`
	CategoryID      int    `json:"category_id"`
	Image           string `json:"image"`
	ProductName     string `json:"name"`
	Size            string `json:"size"`
	Stock           int    `json:"stock"`
	Price           uint64 `json:"price"`
	DiscountedPrice uint64 `json:"discounted_price"`
}

type Banner struct {
	CategoryID         int
	Name               string
	DiscountPercentage int
	Images             []string `gorm:"-"`
}
