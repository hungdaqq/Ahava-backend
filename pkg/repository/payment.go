package repository

import (
	"ahava/pkg/domain"
	"ahava/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type paymentRepository struct {
	DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *paymentRepository {
	return &paymentRepository{
		DB: db,
	}
}

type PaymentRepository interface {
	CreateSePayQR(qr models.CreateQR, user_id uint) error
	SePayWebhook(transaction models.Transaction) error
}

func (p *paymentRepository) CreateSePayQR(qr models.CreateQR, user_id uint) error {
	if err := p.DB.Exec(`INSERT INTO transactions (user_id,order_id,code) VALUES (?,?,?)`, user_id, qr.OrderID, qr.Description).Error; err != nil {
		return err
	}

	return nil
}

func (p *paymentRepository) SePayWebhook(transaction models.Transaction) error {

	result := p.DB.Model(&domain.Transaction{}).Where("code = ?", transaction.Code).Updates(domain.Transaction{
		Gateway:         transaction.Gateway,
		TransactionDate: transaction.TransactionDate,
		AccountNumber:   transaction.AccountNumber,
		Content:         transaction.Content,
		TransferType:    transaction.TransferType,
		TransferAmount:  transaction.TransferAmount,
		Accumulated:     transaction.Accumulated,
		ReferenceCode:   transaction.ReferenceCode,
		Description:     transaction.Description,
	})

	// Check for errors in the update process
	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were affected (to handle case where no transaction with the given code exists)
	if result.RowsAffected == 0 {
		return errors.New("transaction not found")
	}

	return nil
}
