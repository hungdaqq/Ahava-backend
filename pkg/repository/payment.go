package repository

import (
	"ahava/pkg/domain"
	errors "ahava/pkg/utils/errors"
	"ahava/pkg/utils/models"

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
	CreateQR(qr models.CreateQR, user_id uint) error
	SaveTransaction(transaction models.Transaction) (models.Transaction, error)
}

func (p *paymentRepository) CreateQR(qr models.CreateQR, user_id uint) error {

	if err := p.DB.Exec(`INSERT INTO transactions (user_id,order_id,code) VALUES (?,?,?)`, user_id, qr.OrderID, qr.Description).Error; err != nil {
		return err
	}

	return nil
}

func (p *paymentRepository) SaveTransaction(saveTransaction models.Transaction) (models.Transaction, error) {

	var transaction models.Transaction

	result := p.DB.Model(&domain.Transaction{}).Where("code = ?", saveTransaction.Code).Updates(domain.Transaction{
		Gateway:         saveTransaction.Gateway,
		TransactionDate: saveTransaction.TransactionDate,
		AccountNumber:   saveTransaction.AccountNumber,
		Content:         saveTransaction.Content,
		TransferType:    saveTransaction.TransferType,
		TransferAmount:  saveTransaction.TransferAmount,
		Accumulated:     saveTransaction.Accumulated,
		ReferenceCode:   saveTransaction.ReferenceCode,
		Description:     saveTransaction.Description,
	}).Scan(&transaction)

	// Check for errors in the update process
	if result.Error != nil {
		return models.Transaction{}, result.Error
	}

	// Check if any rows were affected (to handle case where no transaction with the given code exists)
	if result.RowsAffected == 0 {
		return models.Transaction{}, errors.ErrEntityNotFound
	}

	return transaction, nil
}
