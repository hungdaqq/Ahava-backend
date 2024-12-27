package models

import "time"

type AdminLogin struct {
	Email    string `json:"email,omitempty" validate:"required"`
	Password string `json:"password" validate:"min=8,max=20"`
}

type AdminDetailsResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name" `
	Email string `json:"email" `
}

type Category struct {
	ID           uint   `json:"id"`
	CategoryName string `json:"category_name" validate:"required"`
	Description  string `json:"description" validate:"required"`
}

type CategoryResponse struct {
	ID          uint    `json:"id" gorm:"unique;not null"`
	Category    string  `json:"Category"`
	ProductName string  `json:"product_name"`
	Size        string  `json:"size"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
}

type UpdateCategory struct {
	CategoryName string `json:"category_name"`
	Description  string `json:"description"`
}

type ListProducts struct {
	Total    int64      `json:"total"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
	Products []Products `json:"products"`
}

type Products struct {
	ID          uint    `json:"id"`
	CategoryID  int     `json:"category_id"`
	ProductName string  `json:"product_name"`
	Size        string  `json:"size"`
	Stock       int     `json:"stock"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	HowToUse    string  `json:"how_to_use"`
}

type ChangePassword struct {
	Oldpassword string `json:"old_password"`
	Password    string `json:"password"`
	Repassword  string `json:"re_password"`
}

type UserSignInResponse struct {
	Id       uint   `json:"id"`
	UserID   uint   `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email" validate:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UserDetailsResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	BirthDate time.Time `json:"birth_date"`
	CreateAt  time.Time `json:"create_at"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password" validate:"required"`
}

type TokenUsers struct {
	Users UserDetailsResponse
	Token string
}

type UserDetails struct {
	Name            string    `json:"name"`
	Email           string    `json:"email" validate:"email"`
	Phone           string    `json:"phone"`
	Password        string    `json:"password"`
	ConfirmPassword string    `json:"confirmpassword"`
	BirthDate       time.Time `json:"birth_date"`
	Address         Address   `json:"address"`
}

type UserDetailsAtAdmin struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Blocked bool   `json:"blocked"`
}

type Search struct {
	Key string `json:"searchkey"`
}

type EditProfile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type CartItem struct {
	CartID              int     `json:"cart_id"`
	ProductID           int     `json:"product_id"`
	ProductName         string  `json:"product_name"`
	Image               string  `json:"image"`
	Quantity            int     `json:"quantity"`
	ItemPrice           float64 `json:"item_price"`
	ItemDiscountedPrice float64 `json:"item_discounted_price"`
}

type UpdateCartItem struct {
	ProductID int  `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type CartDetails struct {
	ID        int `json:"cart_id"`
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
	// CreateAt  time.Time `json:"create_at"`
	// UpdateAt  time.Time `json:"update_at"`
}

type SearchHistory struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	SearchKey string    `json:"search_key"`
	CreateAt  time.Time `json:"create_at"`
}

type CartCheckout struct {
	CartIDs []int `json:"cart_ids"`
}

type CheckOut struct {
	CartItems            []CartItem `json:"cart_items"`
	TotalPrice           float64    `json:"total_price"`
	TotalDiscountedPrice float64    `json:"total_discounted_price"`
}

type Address struct {
	Id       int       `json:"id"`
	UserID   int       `json:"user_id"`
	Name     string    `json:"name"`
	Street   string    `json:"street"`
	Ward     string    `json:"ward"`
	District string    `json:"district"`
	City     string    `json:"city"`
	Phone    string    `json:"phone"`
	Default  bool      `json:"default"`
	Type     string    `json:"type"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

type Wishlist struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
}

type Order struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	AddressID       int       `json:"address_id"`
	PaymentMethodID int       `json:"payment_method_id"`
	FinalPrice      float64   `json:"final_price"`
	CouponUsed      string    `json:"coupon_used"`
	OrderStatus     string    `json:"order_status"`
	CreateAt        time.Time `json:"create_at"`
	UpdateAt        time.Time `json:"update_at"`
}

type PlaceOrder struct {
	UserID          int      `json:"user_id"`
	AddressID       int      `json:"address_id"`
	PaymentMethodID int      `json:"paymentmethod_id"`
	CartCheckOut    CheckOut `json:"cart_checkout"`
	CouponUsed      string   `json:"coupon_used"`
	FinalPrice      float64  `json:"final_price"`
}

type OrderItem struct {
	OrderID         int `json:"order_id"`
	ProductID       int `json:"product_id"`
	Quantity        int `json:"quantity"`
	ItemPrice       int `json:"item_price"`
	DiscountedPrice int `json:"item_discounted_price"`
}

type CreateQR struct {
	OrderID       int     `json:"order_id"`
	AccountNumber string  `json:"account_number"`
	BankName      string  `json:"bank_name"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	Link          string  `json:"link"`
}

type Transaction struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	OrderID         int       `json:"order_id"`
	Gateway         string    `json:"gateway"`
	TransactionDate time.Time `json:"transactionDate"`
	AccountNumber   string    `json:"accountNumber"`
	Code            string    `json:"code"`
	Content         string    `json:"content"`
	TransferType    string    `json:"transferType"`
	TransferAmount  float64   `json:"transferAmount"`
	Accumulated     float64   `json:"accumulated"`
	ReferenceCode   string    `json:"referenceCode"`
	Description     string    `json:"description"`
}
