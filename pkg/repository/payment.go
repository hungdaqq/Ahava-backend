package repository

import (
	"ahava/pkg/domain"
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreateQR(qr models.CreateQR, user_id uint) error
	SaveTransaction(transaction models.Transaction) (models.Transaction, error)
}

type paymentRepository struct {
	DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{
		DB: db,
	}
}

func (r *paymentRepository) CreateQR(qr models.CreateQR, user_id uint) error {

	if err := r.DB.Exec(`INSERT INTO transactions (user_id,order_id,code) VALUES (?,?,?)`, user_id, qr.OrderID, qr.Description).Error; err != nil {
		return err
	}

	return nil
}

func (r *paymentRepository) SaveTransaction(t models.Transaction) (models.Transaction, error) {

	var transaction models.Transaction

	result := r.DB.
		Model(&domain.Transaction{}).
		Where("code = ?", t.Code).
		Updates(domain.Transaction{
			Gateway:         t.Gateway,
			TransactionDate: t.TransactionDate,
			AccountNumber:   t.AccountNumber,
			Content:         t.Content,
			TransferType:    t.TransferType,
			TransferAmount:  t.TransferAmount,
			Accumulated:     t.Accumulated,
			ReferenceCode:   t.ReferenceCode,
			Description:     t.Description,
		}).
		Scan(&transaction)

	// Check for errors in the update process
	if result.Error != nil {
		return models.Transaction{}, result.Error
	}

	// Check if any rows were affected (to handle case where no transaction with the given code exists)
	if result.RowsAffected == 0 {
		return models.Transaction{}, models.ErrEntityNotFound
	}

	return transaction, nil
}
