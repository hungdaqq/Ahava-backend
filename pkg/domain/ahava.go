package domain

import (
	"ahava/pkg/utils/models"
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID       uint   `json:"id" gorm:"unique;not null"`
	Name     string `json:"name" gorm:"validate:required"`
	Email    string `json:"email" gorm:"validate:required"`
	Password string `json:"password" gorm:"validate:required"`
}

type TokenAdmin struct {
	Admin        models.AdminDetailsResponse
	AccessToken  string
	RefreshToken string
}

// type Cart struct {
// 	ID     uint  `json:"id" gorm:"primarykey"`
// 	UserID uint  `json:"user_id" gorm:"not null"`
// 	Users  Users `json:"-" gorm:"foreignkey:UserID"`
// }

type CartItems struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Users     Users     `json:"-" gorm:"foreignkey:UserID"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	Products  Products  `json:"-" gorm:"foreignkey:ProductID;constraint:OnDelete:CASCADE"`
	Quantity  int       `json:"quantity" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime "`
}

type Coupons struct {
	gorm.Model
	Coupon       string `json:"coupon" gorm:"unique;not null"`
	DiscountRate int    `json:"discount_rate" gorm:"not null"`
	Valid        bool   `json:"valid" gorm:"default:true"`
}

type Offer struct {
	ID           int      `json:"id" gorm:"unique;not null"`
	CategoryID   int      `json:"category_id"`
	Category     Category `json:"-" gorm:"foreignkey:CategoryID;constraint:OnDelete:CASCADE"`
	DiscountRate int      `json:"discount_rate"`
	Valid        bool     `gorm:"default:True"`
}

type PaymentMethod struct {
	ID           uint   `gorm:"primarykey"`
	Payment_Name string `json:"payment_name"`
	IsDeleted    bool   `json:"is_deleted" gorm:"default:false"`
}

type Order struct {
	gorm.Model
	UserID          uint          `json:"user_id" gorm:"not null"`
	Users           Users         `json:"-" gorm:"foreignkey:UserID"`
	AddressID       uint          `json:"address_id" gorm:"not null"`
	Address         Address       `json:"-" gorm:"foreignkey:AddressID"`
	PaymentMethodID uint          `json:"paymentmethod_id"`
	PaymentMethod   PaymentMethod `json:"-" gorm:"foreignkey:PaymentMethodID"`
	CouponUsed      string        `json:"coupon_used" gorm:"default:null"`
	FinalPrice      float64       `json:"price"`
	OrderStatus     string        `json:"order_status" gorm:"order_status:4;default:'PENDING';check:order_status IN ('PENDING', 'SHIPPED','DELIVERED','CANCELED','RETURNED')"`
	PaymentStatus   string        `json:"payment_status" gorm:"payment_status:2;default:'NOT PAID';check:payment_status IN ('PAID', 'NOT PAID')"`
}

type OrderItem struct {
	ID         uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID    uint     `json:"order_id"`
	Order      Order    `json:"-" gorm:"foreignkey:OrderID;constraint:OnDelete:CASCADE"`
	ProductID  uint     `json:"product_id"`
	Products   Products `json:"-" gorm:"foreignkey:ProductID"`
	Quantity   int      `json:"quantity"`
	TotalPrice float64  `json:"total_price"`
}

type AdminOrdersResponse struct {
	Pending   []OrderDetails
	Shipped   []OrderDetails
	Delivered []OrderDetails
	Canceled  []OrderDetails
	Returned  []OrderDetails
}

type OrderDetails struct {
	ID            int     `json:"id" gorm:"id"`
	Username      string  `json:"name"`
	Address       string  `json:"address"`
	PaymentMethod string  `json:"payment_method" gorm:"payment_method"`
	Total         float64 `json:"total"`
}

type OrderDetailsWithImages struct {
	OrderDetails  Order
	Images        []string
	PaymentMethod string
}

type Products struct {
	ID          uint     `json:"id" gorm:"unique;not null"`
	CategoryID  int      `json:"category_id"`
	Category    Category `json:"-" gorm:"foreignkey:CategoryID;constraint:OnDelete:CASCADE"`
	ProductName string   `json:"product_name"`
	Image       string   `json:"image"`
	Size        string   `json:"size"`
	Stock       int      `json:"stock"`
	Price       float64  `json:"price"`
	Description string   `json:"description"`
	HowToUse    string   `json:"how_to_use"`
}

type Category struct {
	ID           uint   `json:"id" gorm:"unique;not null"`
	CategoryName string `json:"category_name" gorm:"unique;not null" `
	Description  string `json:"description"`
}

type Users struct {
	ID        uint      `json:"id" gorm:"unique;not null"`
	Name      string    `json:"name"`
	Email     string    `json:"email" validate:"email"`
	Password  string    `json:"password" validate:"min=8,max=20"`
	Phone     string    `json:"phone"`
	BirthDate time.Time `json:"birth_date"`
	// Address      string    `json:"address"`
	Blocked      bool      `json:"blocked" gorm:"default:false"`
	IsAdmin      bool      `json:"is_admin" gorm:"default:false"`
	ReferralCode string    `json:"referral_code"`
	CreateAt     time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateAt     time.Time `json:"update_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

type Address struct {
	Id       uint      `json:"id" gorm:"unique;not null"`
	UserID   uint      `json:"user_id"`
	Users    Users     `json:"-" gorm:"foreignkey:UserID"`
	Name     string    `json:"name" validate:"required"`
	Street   string    `json:"street" validate:"required"`
	Ward     string    `json:"ward" validate:"required"`
	District string    `json:"district" validate:"required"`
	City     string    `json:"city" validate:"required"`
	Phone    string    `json:"phone" gorm:"phone"`
	Default  bool      `json:"default" gorm:"default:false"`
	Type     string    `json:"type" gorm:"default:'HOME';check:type IN ('HOME', 'WORK')"`
	CreateAt time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateAt time.Time `json:"update_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

type Wallet struct {
	ID     int     `json:"id"  gorm:"unique;not null"`
	UserID int     `json:"user_id"`
	Users  Users   `json:"-" gorm:"foreignkey:UserID"`
	Amount float64 `json:"amount" gorm:"default:0"`
}

type Wishlist struct {
	ID        uint     `json:"id" gorm:"primarykey"`
	UserID    uint     `json:"user_id" gorm:"not null"`
	Users     Users    `json:"-" gorm:"foreignkey:UserID"`
	ProductID uint     `json:"product_id" gorm:"not null"`
	Products  Products `json:"-" gorm:"foreignkey:ProductID"`
	// IsDeleted bool     `json:"is_deleted" gorm:"default:false"`
}

type SearchHistory struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Users     Users     `json:"-" gorm:"foreignkey:UserID"`
	SearchKey string    `json:"search_key" gorm:"not null"`
	CreateAt  time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
}
