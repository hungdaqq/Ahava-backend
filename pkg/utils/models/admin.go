package models

type NewPaymentMethod struct {
	PaymentMethod string `json:"payment_method"`
}

type Coupons struct {
	Coupon       string `json:"coupon" gorm:"unique;not null"`
	DiscountRate int    `json:"discount_rate" gorm:"not null"`
	Valid        bool   `json:"valid" gorm:"default:true"`
}
