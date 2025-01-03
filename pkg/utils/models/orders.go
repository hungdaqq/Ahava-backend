package models

type OrderDetails struct {
	ID            int    `json:"order_id"`
	UserName      string `json:"name"`
	AddressID     int    `json:"address_id"`
	PaymentMethod string `json:"payment_method"`
	Amount        uint64 `json:"amount"`
}

type CombinedOrderDetails struct {
	OrderId        string `json:"order_id"`
	FinalPrice     uint64 `json:"final_price"`
	ShipmentStatus string `json:"shipment_status"`
	PaymentStatus  string `json:"payment_status"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	HouseName      string `json:"house_name" validate:"required"`
	State          string `json:"state" validate:"required"`
	Pin            string `json:"pin" validate:"required"`
	Street         string `json:"street"`
	Province       string `json:"province"`
}

type OrderPaymentDetails struct {
	UserID     int    `json:"user_id"`
	Username   string `json:"username"`
	Razor_id   string `josn:"razor_id"`
	OrderID    int    `json:"order_id"`
	FinalPrice uint64 `json:"final_price"`
}

type EditOrderStatus struct {
	Orderid uint   `json:"order_id"`
	Status  string `json:"order_status"`
}

type IndividualOrderDetails struct {
	OrderID       int
	Address       string
	Phone         string
	Products      []Products `gorm:"-"`
	TotalAmount   uint64
	CouponUsed    string
	OrderStatus   string
	PaymentStatus string
}

// type ProductDetails struct {
// 	ProductName string
// 	Image       string
// 	Quantity    int
// 	Amount      float64
// }
