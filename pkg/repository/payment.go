package repository

import (
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
	CreateSePayQR(qr models.CreateQR, user_id int) error
	SePayWebhook(transaction models.Transaction) error
}

func (p *paymentRepository) CreateSePayQR(qr models.CreateQR, user_id int) error {
	if err := p.DB.Exec(`INSERT INTO transactions (user_id,order_id,code) VALUES (?,?,?)`, user_id, qr.OrderID, qr.Description).Error; err != nil {
		return err
	}

	return nil
}

func (p *paymentRepository) SePayWebhook(transaction models.Transaction) error {

	err := p.DB.Exec(`UPDATE transactions SET gateway=$1,transaction_date=$2,account_number=$3,content=$4,transfer_type=$5,
					transfer_amount=$6,accumulated=$7,reference_code=$8,description=$9 WHERE code=$10`,
		transaction.Gateway, transaction.TransactionDate, transaction.AccountNumber, transaction.Content, transaction.TransferType,
		transaction.TransferAmount, transaction.Accumulated, transaction.ReferenceCode, transaction.Description, transaction.Code).Error
	if err != nil {
		return err
	}

	return nil
}
