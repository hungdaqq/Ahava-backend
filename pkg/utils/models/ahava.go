package models

import (
	"time"

	"github.com/lib/pq"
)

type AdminLogin struct {
	Email    string `json:"email,omitempty" validate:"required"`
	Password string `json:"password" validate:"min=8,max=20"`
}

type AdminDetailsResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name" `
	Email string `json:"email" `
}

type Category struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateCategory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ListProducts struct {
	Total    int64     `json:"total"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
	Products []Product `json:"products"`
}

type Product struct {
	ID               uint           `json:"id"`
	CategoryID       uint           `json:"category_id"`
	Name             string         `json:"name"`
	Size             string         `json:"size"`
	Stock            uint           `json:"stock"`
	DefaultImage     string         `json:"default_image"`
	Images           pq.StringArray `json:"images"`
	Price            uint64         `json:"price"`
	DiscountedPrice  uint64         `json:"discounted_price"`
	ShortDescription string         `json:"short_description"`
	Description      string         `json:"description"`
	HowToUse         string         `json:"how_to_use"`
	IsFeatured       bool           `json:"is_featured"`
}

type WishlistProduct struct {
	ID              uint   `json:"id"`
	ProductID       uint   `json:"product_id"`
	Name            string `json:"name"`
	Size            string `json:"size"`
	Stock           uint   `json:"stock"`
	DefaultImage    string `json:"default_image"`
	Price           uint64 `json:"price"`
	DiscountedPrice uint64 `json:"discounted_price"`
}

type ChangePassword struct {
	Oldpassword string `json:"old_password"`
	Password    string `json:"password"`
	Repassword  string `json:"re_password"`
}

type UserSignInResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
	BirthDate time.Time `json:"birth_date"`
	Gender    string    `json:"gender"`
}

type UserDetailsResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenUsers struct {
	Users UserDetailsResponse
	Token string
}

type UserDetails struct {
	Name            string    `json:"name"`
	Username        string    `json:"username"`
	Email           string    `json:"email" validate:"email"`
	Gender          string    `json:"gender"`
	Phone           string    `json:"phone"`
	Password        string    `json:"password"`
	ConfirmPassword string    `json:"confirmpassword"`
	BirthDate       time.Time `json:"birth_date"`
	Address         Address   `json:"address"`
}

type UserDetailsAtAdmin struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Blocked bool   `json:"blocked"`
}

type Search struct {
	Key string `json:"searchkey"`
}

type EditProfile struct {
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	BirthDate time.Time `json:"birth_date"`
	Gender    string    `json:"gender"`
}

type CartItem struct {
	CartID              uint   `json:"cart_id"`
	ProductID           uint   `json:"product_id"`
	Name                string `json:"name"`
	DefaultImage        string `json:"default_image"`
	Quantity            uint   `json:"quantity"`
	ItemPrice           uint64 `json:"item_price"`
	ItemDiscountedPrice uint64 `json:"item_discounted_price"`
}

type UpdateCartItem struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type CartDetails struct {
	ID        uint `json:"cart_id"`
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
	// CreateAt  time.Time `json:"create_at"`
	// UpdateAt  time.Time `json:"update_at"`
}

type SearchHistory struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	SearchKey string    `json:"search_key"`
	CreateAt  time.Time `json:"create_at"`
}

type CartCheckout struct {
	CartIDs []uint `json:"cart_ids"`
}

type CheckOut struct {
	CartItems            []CartItem `json:"cart_items"`
	TotalPrice           uint64     `json:"total_price"`
	TotalDiscountedPrice uint64     `json:"total_discounted_price"`
}

type Address struct {
	ID       uint      `json:"id"`
	UserID   uint      `json:"user_id"`
	Name     string    `json:"name"`
	Street   string    `json:"street"`
	Ward     string    `json:"ward"`
	District string    `json:"district"`
	Province string    `json:"province"`
	Phone    string    `json:"phone"`
	Default  bool      `json:"default"`
	Type     string    `json:"type"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

type Wishlist struct {
	ID        uint `json:"id"`
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
}

type AddToWishlist struct {
	ProductID uint `json:"product_id"`
}

type Order struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	Address       string    `json:"address"`
	PaymentMethod string    `json:"payment_method"`
	FinalPrice    uint64    `json:"final_price"`
	Coupon        string    `json:"coupon"`
	OrderStatus   string    `json:"order_status"`
	PaymentStatus string    `json:"payment_status"`
	CreateAt      time.Time `json:"create_at"`
	UpdateAt      time.Time `json:"update_at"`
}

type PlaceOrder struct {
	UserID        uint   `json:"user_id"`
	Address       string `json:"address"`
	PaymentMethod string `json:"payment_method"`
	CartIDs       []uint `json:"cart_ids"`
	Coupon        string `json:"coupon"`
}

type OrderItem struct {
	OrderID         uint   `json:"order_id"`
	ProductID       uint   `json:"product_id"`
	Quantity        uint   `json:"quantity"`
	ItemPrice       uint64 `json:"item_price"`
	DiscountedPrice uint64 `json:"item_discounted_price"`
}

type CreateQR struct {
	OrderID       uint   `json:"order_id" validate:"required"`
	AccountNumber string `json:"account_number"`
	BankName      string `json:"bank_name"`
	Amount        uint64 `json:"amount" validate:"required"`
	Description   string `json:"description"`
	Link          string `json:"link"`
}

type Transaction struct {
	ID              uint   `json:"id"`
	UserID          uint   `json:"user_id"`
	OrderID         uint   `json:"order_id"`
	Gateway         string `json:"gateway"`
	TransactionDate string `json:"transactionDate"`
	AccountNumber   string `json:"accountNumber"`
	Code            string `json:"code"`
	Content         string `json:"content"`
	TransferType    string `json:"transferType"`
	TransferAmount  uint64 `json:"transferAmount"`
	Accumulated     uint64 `json:"accumulated"`
	ReferenceCode   string `json:"referenceCode"`
	Description     string `json:"description"`
}

type Offer struct {
	ProductID uint      `json:"product_id"`
	OfferRate uint      `json:"offer_rate"`
	ExpireAt  time.Time `json:"expire_at"`
}

type Province struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	TypeText string `json:"typeText"`
	Slug     string `json:"slug"`
}

type Provinces struct {
	Total     int        `json:"total"`
	Provinces []Province `json:"data"`
}
