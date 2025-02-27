package models

import (
	"errors"
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

type Price struct {
	ID            uint   `json:"id"`
	Size          string `json:"size"`
	Image         string `json:"image"`
	OriginalPrice uint64 `json:"original_price" gorm:"default:1"`
	DiscountPrice uint64 `json:"discount_price"`
}

type Product struct {
	ID               uint           `json:"id"`
	Category         string         `json:"category"`
	Name             string         `json:"name"`
	Code             string         `json:"code"`
	DefaultImage     string         `json:"default_image"`
	Images           pq.StringArray `json:"images"`
	Stock            uint           `json:"stock"`
	Type             string         `json:"type"`
	Tag              string         `json:"tag"`
	Price            []Price        `json:"price"`
	ShortDescription string         `json:"short_description"`
	Description      string         `json:"description"`
	HowToUse         string         `json:"how_to_use"`
	IsFeatured       bool           `json:"is_featured"`
}

type WishlistProduct struct {
	ID            uint   `json:"id"`
	ProductID     uint   `json:"product_id"`
	Name          string `json:"name"`
	DefaultImage  string `json:"default_image"`
	Size          string `json:"size"`
	OriginalPrice uint   `json:"original_price"`
	DiscountPrice uint   `json:"discount_price"`
}

type ChangePassword struct {
	Oldpassword string `json:"old_password" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Repassword  string `json:"re_password" validate:"required"`
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
	Username string `json:"username" `
	Email    string `json:"email" `
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
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender"`
	IsBlocked bool   `json:"is_blocked"`
}

type ListUsers struct {
	Total  int64                `json:"total"`
	Limit  int                  `json:"limit"`
	Offset int                  `json:"offset"`
	Users  []UserDetailsAtAdmin `json:"users"`
}

type Search struct {
	Key string `json:"searchkey" validate:"required"`
}

type EditProfile struct {
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	BirthDate time.Time `json:"birth_date"`
	Gender    string    `json:"gender"`
}

type CartItem struct {
	ID                uint   `json:"id"`
	ProductID         uint   `json:"product_id"`
	Name              string `json:"name"`
	DefaultImage      string `json:"default_image"`
	Size              string `json:"size"`
	Quantity          uint   `json:"quantity"`
	OriginalPrice     uint64 `json:"original_price"`
	DiscountPrice     uint64 `json:"discount_price"`
	ItemPrice         uint64 `json:"item_price"`
	ItemDiscountPrice uint64 `json:"item_discount_price"`
}

type UpdateCartItem struct {
	ProductID uint   `json:"product_id" validate:"required"`
	Quantity  uint   `json:"quantity" validate:"required"`
	Size      string `json:"size" validate:"required"`
}

type CartDetails struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id"`
	ProductID     uint   `json:"product_id"`
	Size          string `json:"size"`
	Quantity      uint   `json:"quantity"`
	OriginalPrice uint64 `json:"original_price"`
	DiscountPrice uint64 `json:"discount_price"`
}

type SearchHistory struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	SearchKey string    `json:"search_key"`
	CreateAt  time.Time `json:"created_at"`
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
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	Name         string `json:"name"`
	Street       string `json:"street"`
	Ward         string `json:"ward"`
	WardCode     string `json:"ward_code"`
	District     string `json:"district"`
	DistrictCode string `json:"district_code"`
	Province     string `json:"province"`
	ProvinceCode string `json:"province_code"`
	Phone        string `json:"phone"`
	Default      bool   `json:"default"`
	Type         string `json:"type"`
}

type Wishlist struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	ProductID uint   `json:"product_id"`
	Size      string `json:"size"`
}

type AddToWishlist struct {
	ProductID uint   `json:"product_id" validate:"required"`
	Size      string `json:"size" validate:"required"`
}

type ListOrders struct {
	Total  int64          `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
	Orders []OrderDetails `json:"orders"`
}

type OrderDetails struct {
	Order
	Details []OrderItem `json:"details"`
}

type Order struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	PaymentMethod string `json:"payment_method"`
	FinalPrice    uint64 `json:"final_price"`
	Coupon        string `json:"coupon"`
	OrderStatus   string `json:"order_status"`
	PaymentStatus string `json:"payment_status"`
}

type PlaceOrder struct {
	UserID        uint   `json:"user_id"`
	Address       string `json:"address"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	PaymentMethod string `json:"payment_method"`
	CartIDs       []uint `json:"cart_ids"`
	Coupon        string `json:"coupon"`
}

type OrderItem struct {
	OrderID             uint   `json:"order_id"`
	ProductID           uint   `json:"product_id"`
	Size                string `json:"size"`
	Quantity            uint   `json:"quantity"`
	OriginalPrice       uint64 `json:"original_price"`
	DiscountPrice       uint64 `json:"discounted_price"`
	ItemPrice           uint64 `json:"item_price"`
	ItemDiscountedPrice uint64 `json:"item_discount_price"`
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

type News struct {
	ID             uint   `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Content        string `json:"content"`
	DefaultImage   string `json:"default_image"`
	IsFeatured     bool   `json:"is_featured"`
	IsHomepage     bool   `json:"is_homepage"`
	IsDisplay      bool   `json:"is_display"`
	Category       string `json:"category"`
	TitleSEO       string `json:"title_seo"`
	DescriptionSEO string `json:"description_seo"`
	LinkSEO        string `json:"link_seo"`
}

type ListNews struct {
	Total  int64  `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	News   []News `json:"news"`
}

var (
	ErrEntityNotFound  = errors.New("entity not found")
	ErrInternalServer  = errors.New("internal server error")
	ErrBadRequest      = errors.New("bad request")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrConflict        = errors.New("conflict")
	ErrInvalidToken    = errors.New("invalid token")
	ErrCreateToken     = errors.New("error in creating token")
	ErrValidateOTP     = errors.New("failed to validate otp")
	ErrAlreadyExists   = errors.New("entity already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrMalformedEntity = errors.New("malformed entiry")
)
