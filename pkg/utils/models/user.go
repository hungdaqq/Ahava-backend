package models

// import "time"

// type UserDetails struct {
// 	Name            string    `json:"name"`
// 	Email           string    `json:"email" validate:"email"`
// 	Phone           string    `json:"phone"`
// 	Password        string    `json:"password"`
// 	ConfirmPassword string    `json:"confirmpassword"`
// 	BirthDate       time.Time `json:"birth_date"`
// 	Address         string    `json:"address"`
// }

type Address struct {
	Id        uint   `json:"id" gorm:"unique;not null"`
	UserID    uint   `json:"user_id"`
	Name      string `json:"name" validate:"required"`
	HouseName string `json:"house_name" validate:"required"`
	Street    string `json:"street" validate:"required"`
	City      string `json:"city" validate:"required"`
	State     string `json:"state" validate:"required"`
	Pin       string `json:"pin" validate:"required"`
}

type AddAddress struct {
	Name      string `json:"name" validate:"required"`
	HouseName string `json:"house_name" validate:"required"`
	Street    string `json:"street" validate:"required"`
	City      string `json:"city" validate:"required"`
	State     string `json:"state" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Pin       string `json:"pin" validate:"required"`
}

// type ForgotPasswordSend struct {
// 	Phone string `json:"phone"`
// }

// type ForgotVerify struct {
// 	Phone       string `json:"phone"`
// 	Otp         string `json:"otp"`
// 	NewPassword string `json:"newpassword"`
// }


// type CheckOut struct {
// 	CartID          int
// 	Addresses       []Address
// 	Products        []Cart
// 	PaymentMethods  []PaymentMethod
// 	TotalPrice      float64
// 	DiscountedPrice float64
// }
