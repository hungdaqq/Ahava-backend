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

type CartItem struct {
	gorm.Model
	ID        uint    `json:"id" gorm:"primarykey"`
	UserID    uint    `json:"user_id" gorm:"not null"`
	User      User    `json:"-" gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"-" gorm:"foreignkey:ProductID;constraint:OnDelete:CASCADE"`
	Size      string  `json:"size" gorm:"not null"`
	Quantity  uint    `json:"quantity" gorm:"default:1;check:quantity>0"`
}

type Coupons struct {
	gorm.Model
	Coupon       string `json:"coupon" gorm:"unique;not null"`
	DiscountRate int    `json:"discount_rate" gorm:"not null"`
	Valid        bool   `json:"valid" gorm:"default:true"`
}

type PaymentMethod struct {
	ID          uint   `json:"id" gorm:"primarykey"`
	PaymentName string `json:"payment_name"`
	Enable      bool   `json:"enable" gorm:"default:true"`
}

type Order struct {
	gorm.Model
	UserID        uint   `json:"user_id" gorm:"not null"`
	User          User   `json:"-" gorm:"foreignkey:UserID"`
	Name          string `json:"name" gorm:"not null"`
	Phone         string `json:"phone" gorm:"not null"`
	Address       string `json:"address" gorm:"not null"`
	PaymentMethod string `json:"payment_method"`
	Coupon        string `json:"coupon" gorm:"default:null"`
	FinalPrice    uint64 `json:"price" gorm:"not null"`
	OrderStatus   string `json:"order_status" gorm:"order_status:10;default:'UNCONFIRMED';check:order_status IN ('UNCONFIRMED', 'PREPARING','SHIPPING','DELIVERED','CANCELED','RETURNED')"`
	PaymentStatus string `json:"payment_status" gorm:"payment_status:2;default:'NOT PAID';check:payment_status IN ('PAID', 'NOT PAID', 'INCOMPLETE')"`
}

type OrderItem struct {
	gorm.Model
	OrderID           uint    `json:"order_id" gorm:"not null"`
	Order             Order   `json:"-" gorm:"foreignkey:OrderID;constraint:OnDelete:CASCADE"`
	ProductID         uint    `json:"product_id" gorm:"not null"`
	Product           Product `json:"-" gorm:"foreignkey:ProductID"`
	Size              string  `json:"size" gorm:"not null"`
	Quantity          uint    `json:"quantity" gorm:"not null"`
	OriginalPrice     uint64  `json:"original_price" gorm:"not null"`
	DiscountPrice     uint64  `json:"discounted_price" gorm:"not null"`
	ItemPrice         uint64  `json:"item_price" gorm:"not null"`
	ItemDiscountPrice uint64  `json:"item_discount_price" gorm:"not null"`
}

type OrderDetails struct {
	gorm.Model
	Username      string `json:"name"`
	Address       string `json:"address"`
	PaymentMethod string `json:"payment_method" gorm:"payment_method"`
	Total         uint64 `json:"total"`
}

type Price struct {
	gorm.Model
	ProductID     uint    `json:"product_id"`
	Product       Product `json:"-" gorm:"foreignkey:ProductID;constraint:OnDelete:CASCADE"`
	Size          string  `json:"size"`
	Image         string  `json:"image" gorm:"default:'https://minio.ahava.com.vn/ahava/default_product_image.png'"`
	OriginalPrice uint64  `json:"original_price" gorm:"default:1"`
	DiscountPrice uint64  `json:"discount_price"`
}

type Product struct {
	gorm.Model
	Category         string         `json:"category" gorm:"not null"`
	Name             string         `json:"name" gorm:"default:Ahava Product"`
	Code             string         `json:"code" gorm:"default:AVAHA"`
	DefaultImage     string         `json:"default_image" gorm:"default:'https://minio.ahava.com.vn/ahava/default_product_image.png'"`
	Images           pq.StringArray `json:"images" gorm:"type:varchar[]"`
	Stock            uint           `json:"stock" gorm:"default:100"`
	Type             string         `json:"type"`
	Tag              string         `json:"tag"`
	ShortDescription string         `json:"short_description"`
	Description      string         `json:"description"`
	HowToUse         string         `json:"how_to_use"`
	IsFeatured       *bool          `json:"is_featured" gorm:"default:false"`
	IsHidden         *bool          `json:"is_hidden" gorm:"default:false"`
}

type User struct {
	gorm.Model
	Name         string    `json:"name" gorm:"not null"`
	Username     string    `json:"username" gorm:"unique;not null"`
	Email        string    `json:"email" validate:"email" gorm:"unique;not null"`
	Gender       string    `json:"gender" gorm:"gender:2;default:'MALE';check:gender IN ('MALE', 'FEMALE', 'OTHER')"`
	Password     string    `json:"password" validate:"min=8,max=20"`
	Phone        string    `json:"phone"`
	BirthDate    time.Time `json:"birth_date"`
	IsBlocked    bool      `json:"is_blocked" gorm:"default:false"`
	IsAdmin      bool      `json:"is_admin" gorm:"default:false"`
	ReferralCode string    `json:"referral_code"`
}

type Address struct {
	gorm.Model
	UserID       uint   `json:"user_id" gorm:"not null"`
	User         User   `json:"-" gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE"`
	Name         string `json:"name" gorm:"not null"`
	Street       string `json:"street" gorm:"not null"`
	Ward         string `json:"ward" gorm:"not null"`
	WardCode     string `json:"ward_code" gorm:"not null"`
	District     string `json:"district" gorm:"not null"`
	DistrictCode string `json:"district_code" gorm:"not null"`
	Province     string `json:"province" gorm:"not null"`
	ProvinceCode string `json:"province_code" gorm:"not null"`
	Phone        string `json:"phone" gorm:"not null"`
	Default      bool   `json:"default" gorm:"default:false"`
	Type         string `json:"type" gorm:"default:'HOME';check:type IN ('HOME', 'WORK')"`
}

type News struct {
	gorm.Model
	Title          string `json:"title" gorm:"not null"`
	Description    string `json:"description" gorm:"not null"`
	Content        string `json:"content" gorm:"not null"`
	DefaultImage   string `json:"default_image" gorm:"not null"`
	IsFeatured     *bool  `json:"is_featured" gorm:"default:false"`
	IsHomepage     *bool  `json:"is_homepage" gorm:"default:false"`
	IsDisplay      *bool  `json:"is_display" gorm:"default:false"`
	Category       string `json:"category" gorm:"default:'NEWS';check:category IN ('NEWS', 'TIPS', 'NEW_ARRIVALS', 'RECRUITMENT')"`
	TitleSEO       string `json:"title_seo"`
	DescriptionSEO string `json:"description_seo"`
	LinkSEO        string `json:"link_seo"`
}

type Wishlist struct {
	gorm.Model
	UserID    uint    `json:"user_id" gorm:"not null"`
	User      User    `json:"-" gorm:"foreignkey:UserID;constraint:OnDelete:CASCADE"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"-" gorm:"foreignkey:ProductID;constraint:OnDelete:CASCADE"`
	Size      string  `json:"size" gorm:"not null"`
	IsDeleted bool    `json:"is_deleted" gorm:"default:false"`
}

type Transaction struct {
	gorm.Model
	UserID          uint   `json:"user_id" gorm:"not null"`
	User            User   `json:"-" gorm:"foreignkey:UserID"`
	OrderID         uint   `json:"order_id" gorm:"unique;not null"`
	Order           Order  `json:"-" gorm:"foreignkey:OrderID"`
	Gateway         string `json:"gateway"`
	TransactionDate string `json:"transaction_date"`
	AccountNumber   string `json:"account_number"`
	Code            string `json:"code"`
	Content         string `json:"content"`
	TransferType    string `json:"transfer_type"`
	TransferAmount  uint64 `json:"transfer_amount"`
	Accumulated     uint64 `json:"accumulated"`
	ReferenceCode   string `json:"reference_code"`
	Description     string `json:"description"`
}

type RequestTransaction struct {
	gorm.Model
	Method       string `json:"method"`
	Path         string `json:"path"`
	StatusCode   int    `json:"status_code"`
	ClientIP     string `json:"client_ip"`
	Latency      string `json:"latency"`
	BodySize     int    `json:"body_size"`
	ErrorMessage string `json:"error_message"`
}
