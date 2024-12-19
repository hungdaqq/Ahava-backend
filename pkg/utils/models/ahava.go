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

type CategoryProducts struct {
	CategoryID   int        `json:"category_id"`
	CategoryName string     `json:"category_name"`
	Products     []Products `json:"products"`
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
	CartID      int    `json:"cart_id"`
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Image       string `json:"image"`
	// Category_id     int     `json:"category_id"`
	Quantity int `json:"quantity"`
	// StockAvailable  int     `json:"stock"`
	ItemPrice       float64 `json:"item_price"`
	DiscountedPrice float64 `json:"discounted_price"`
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
	Addresses       []Address       `json:"addresses"`
	CartItems       []CartItem      `json:"cart_items"`
	PaymentMethods  []PaymentMethod `json:"payment_methods"`
	TotalPrice      float64         `json:"total_price"`
	DiscountedPrice float64         `json:"discounted_price"`
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
