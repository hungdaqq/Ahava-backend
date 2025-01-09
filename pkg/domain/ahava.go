package domain

import (
	"ahava/pkg/utils/models"
	"time"

	"github.com/lib/pq"

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

type CartItem struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"-" gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	Product   Product   `json:"-" gorm:"foreignkey:ProductID;constraint:OnDelete:CASCADE"`
	Quantity  uint      `json:"quantity" gorm:"default:1;check:quantity > 0"`
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
	Product   Product   `json:"-" gorm:"foreignkey:ProductID;constraint:OnDelete:CASCADE"`
	OfferRate uint      `json:"offer_rate" validate:"min=0,max=100" gorm:"not null"`
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
	User          User      `json:"-" gorm:"foreignkey:UserID"`
	Address       string    `json:"address" gorm:"not null"`
	PaymentMethod string    `json:"payment_method"`
	Coupon        string    `json:"coupon" gorm:"default:null"`
	FinalPrice    uint64    `json:"price" gorm:"not null"`
	OrderStatus   string    `json:"order_status" gorm:"order_status:10;default:'UNCONFIRMED';check:order_status IN ('UNCONFIRMED', 'PREPARING','SHIPPING','DELIVERED','CANCELED','RETURNED')"`
	PaymentStatus string    `json:"payment_status" gorm:"payment_status:2;default:'NOT PAID';check:payment_status IN ('PAID', 'NOT PAID', 'INCOMPLETE')"`
	CreateAt      time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateAt      time.Time `json:"update_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

type OrderItem struct {
	ID                  uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID             uint    `json:"order_id" gorm:"not null"`
	Order               Order   `json:"-" gorm:"foreignkey:OrderID;constraint:OnDelete:CASCADE"`
	ProductID           uint    `json:"product_id" gorm:"not null"`
	Product             Product `json:"-" gorm:"foreignkey:ProductID"`
	Quantity            uint    `json:"quantity" gorm:"not null"`
	ItemPrice           uint64  `json:"item_price" gorm:"not null"`
	ItemDiscountedPrice uint64  `json:"item_discounted_price" gorm:"not null"`
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

type Product struct {
	ID               uint           `json:"id" gorm:"primarykey"`
	CategoryID       uint           `json:"category_id" gorm:"not null"`
	Category         Category       `json:"-" gorm:"foreignkey:CategoryID;constraint:OnDelete:CASCADE"`
	Name             string         `json:"name" gorm:"default:Ahava Product"`
	DefaultImage     string         `json:"default_image" gorm:"not null"`
	Images           pq.StringArray `json:"images" gorm:"type:varchar[]"`
	Size             string         `json:"size"`
	Stock            uint           `json:"stock" gorm:"default:100"`
	Price            uint64         `json:"price" gorm:"default:0"`
	ShortDescription string         `json:"short_description"`
	Description      string         `json:"description"`
	HowToUse         string         `json:"how_to_use"`
	IsFeatured       bool           `json:"is_featured" gorm:"default:false"`
	CreateAt         time.Time      `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateAt         time.Time      `json:"update_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

type Category struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreateAt    time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateAt    time.Time `json:"update_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

type User struct {
	ID           uint      `json:"id" gorm:"unique;not null"`
	Name         string    `json:"name" gorm:"not null"`
	Username     string    `json:"username" gorm:"unique;not null"`
	Email        string    `json:"email" validate:"email" gorm:"unique;not null"`
	Gender       string    `json:"gender" gorm:"gender:2;default:'MALE';check:gender IN ('MALE', 'FEMALE', 'OTHER')"`
	Password     string    `json:"password" validate:"min=8,max=20"`
	Phone        string    `json:"phone" gorm:"not null"`
	BirthDate    time.Time `json:"birth_date"`
	IsBlocked    bool      `json:"is_blocked" gorm:"default:false"`
	IsAdmin      bool      `json:"is_admin" gorm:"default:false"`
	ReferralCode string    `json:"referral_code"`
	CreateAt     time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateAt     time.Time `json:"update_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

type Address struct {
	ID       uint      `json:"id" gorm:"unique;not null"`
	UserID   uint      `json:"user_id" gorm:"not null"`
	User     User      `json:"-" gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE"`
	Name     string    `json:"name" gorm:"not null"`
	Street   string    `json:"street" gorm:"not null"`
	Ward     string    `json:"ward" gorm:"not null"`
	District string    `json:"district" gorm:"not null"`
	Province string    `json:"province" gorm:"not null"`
	Phone    string    `json:"phone" gorm:"not null"`
	Default  bool      `json:"default" gorm:"default:false"`
	Type     string    `json:"type" gorm:"default:'HOME';check:type IN ('HOME', 'WORK')"`
	CreateAt time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateAt time.Time `json:"update_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

type Wallet struct {
	ID     int    `json:"id" gorm:"unique;not null"`
	Userid uint   `json:"user_id"`
	User   User   `json:"-" gorm:"foreignkey:UserID"`
	Amount uint64 `json:"amount" gorm:"default:0"`
}

type Wishlist struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"-" gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	Product   Product   `json:"-" gorm:"foreignkey:ProductID;constraint:OnDelete:CASCADE"`
	IsDeleted bool      `json:"is_deleted" gorm:"default:false"`
	CreateAt  time.Time `json:"create_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdateAt  time.Time `json:"update_at" gorm:"default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

type Transaction struct {
	ID              uint      `json:"id" gorm:"unique;not null"`
	UserID          uint      `json:"user_id" gorm:"not null"`
	User            User      `json:"-" gorm:"foreignkey:UserID"`
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
