package domain

import (
	"ahava/pkg/utils/models"
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID       uint   `json:"id" gorm:"primarykey"`
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
	Users     Users     `json:"-" gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	Products  Products  `json:"-" gorm:"foreignkey:ProductID;constraint:OnDelete:CASCADE"`
	Quantity  uint      `json:"quantity" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateAt  time.Time `json:"update_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

type Coupons struct {
	gorm.Model
	Coupon       string `json:"coupon" gorm:"unique;not null"`
	DiscountRate int    `json:"discount_rate" gorm:"not null"`
	Valid        bool   `json:"valid" gorm:"default:true"`
}

type Offer struct {
	ProductID uint      `json:"product_id" gorm:"not null;unique"`
	Product   Products  `json:"-" gorm:"foreignkey:ProductID;constraint:OnDelete:CASCADE"`
	OfferRate uint      `json:"offer_rate" validate:"min=0,max=100"` // Validation: 0-100
	ExpireAt  time.Time `json:"expire_at"`
	CreateAt  time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateAt  time.Time `json:"update_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

type PaymentMethod struct {
	ID          uint   `json:"id" gorm:"primarykey"`
	PaymentName string `json:"payment_name"`
	Enable      bool   `json:"enable" gorm:"default:true"`
}

type Order struct {
	ID            uint      `json:"id" gorm:"primarykey"`
	UserID        uint      `json:"user_id" gorm:"not null"`
	Users         Users     `json:"-" gorm:"foreignkey:UserID"`
	Address       string    `json:"address"`
	PaymentMethod string    `json:"payment_method"`
	Coupon        string    `json:"coupon" gorm:"default:null"`
	FinalPrice    uint64    `json:"price"`
	OrderStatus   string    `json:"order_status" gorm:"order_status:10;default:'PREPARE';check:order_status IN ('PREPARE','SHIP','DELIVER','CANCEL','RETURN')"`
	PaymentStatus string    `json:"payment_status" gorm:"payment_status:2;default:'NOT PAID';check:payment_status IN ('PAID', 'NOT PAID', 'INCOMPLETE')"`
	CreateAt      time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateAt      time.Time `json:"update_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

type OrderItem struct {
	ID                  uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID             uint     `json:"order_id" gorm:"not null"`
	Order               Order    `json:"-" gorm:"foreignkey:OrderID;constraint:OnDelete:CASCADE"`
	ProductID           uint     `json:"product_id" gorm:"not null"`
	Products            Products `json:"-" gorm:"foreignkey:ProductID"`
	Quantity            uint     `json:"quantity"`
	ItemPrice           uint64   `json:"item_price"`
	ItemDiscountedPrice uint64   `json:"item_discounted_price"`
}

type AdminOrdersResponse struct {
	Pending   []OrderDetails
	Shipped   []OrderDetails
	Delivered []OrderDetails
	Canceled  []OrderDetails
	Returned  []OrderDetails
}

type OrderDetails struct {
	ID            uint   `json:"id" gorm:"id"`
	Username      string `json:"name"`
	Address       string `json:"address"`
	PaymentMethod string `json:"payment_method" gorm:"payment_method"`
	Total         uint64 `json:"total"`
}

type OrderDetailsWithImages struct {
	OrderDetails  Order
	Images        []string
	PaymentMethod string
}

type Products struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	CategoryID  uint      `json:"category_id" gorm:"not null"`
	Category    Category  `json:"-" gorm:"foreignkey:CategoryID;constraint:OnDelete:CASCADE"`
	ProductName string    `json:"product_name"`
	Image       string    `json:"image"`
	Size        string    `json:"size"`
	Stock       uint      `json:"stock"`
	Price       uint64    `json:"price"`
	Description string    `json:"description"`
	HowToUse    string    `json:"how_to_use"`
	IsFeatured  bool      `json:"is_featured" gorm:"default:false"`
	CreateAt    time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateAt    time.Time `json:"update_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

type Category struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	CategoryName string    `json:"category_name"`
	Description  string    `json:"description"`
	CreateAt     time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateAt     time.Time `json:"update_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
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
	ID       uint      `json:"id" gorm:"unique;not null"`
	UserID   uint      `json:"user_id" gorm:"not null"`
	Users    Users     `json:"-" gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE"`
	Name     string    `json:"name" validate:"required"`
	Street   string    `json:"street" validate:"required"`
	Ward     string    `json:"ward" validate:"required"`
	District string    `json:"district" validate:"required"`
	City     string    `json:"city" validate:"required"`
	Phone    string    `json:"phone"`
	Default  bool      `json:"default" gorm:"default:false"`
	Type     string    `json:"type" gorm:"default:'HOME';check:type IN ('HOME', 'WORK')"`
	CreateAt time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateAt time.Time `json:"update_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

type Wallet struct {
	ID     int    `json:"id" gorm:"unique;not null"`
	Userid uint   `json:"user_id"`
	Users  Users  `json:"-" gorm:"foreignkey:UserID"`
	Amount uint64 `json:"amount" gorm:"default:0"`
}

type Wishlist struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Users     Users     `json:"-" gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	Products  Products  `json:"-" gorm:"foreignkey:ProductID;constraint:OnDelete:CASCADE"`
	IsDeleted bool      `json:"is_deleted" gorm:"default:false"`
	CreateAt  time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateAt  time.Time `json:"update_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

type SearchHistory struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Users     Users     `json:"-" gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE"`
	SearchKey string    `json:"search_key" gorm:"not null"`
	CreateAt  time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
}

type Transaction struct {
	ID              uint      `json:"id" gorm:"unique;not null"`
	UserID          uint      `json:"user_id" gorm:"not null"`
	Users           Users     `json:"-" gorm:"foreignkey:UserID"`
	OrderID         uint      `json:"order_id" gorm:"unique;not null"`
	Order           Order     `json:"-" gorm:"foreignkey:OrderID"`
	Gateway         string    `json:"gateway"`
	TransactionDate string    `json:"transaction_date"`
	AccountNumber   string    `json:"account_number"`
	Code            string    `json:"code"`
	Content         string    `json:"content"`
	TransferType    string    `json:"transfer_type"`
	TransferAmount  uint64    `json:"transfer_amount"`
	Accumulated     uint64    `json:"accumulated"`
	ReferenceCode   string    `json:"reference_code"`
	Description     string    `json:"description"`
	CreateAt        time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateAt        time.Time `json:"update_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}
