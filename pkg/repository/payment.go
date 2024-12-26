package repository

import "gorm.io/gorm"

type paymentRepository struct {
	DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *paymentRepository {
	return &paymentRepository{
		DB: db,
	}
}

type PaymentRepository interface {
	// GetPaymentMethods() ([]PaymentMethod, error)
}
