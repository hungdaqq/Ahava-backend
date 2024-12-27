package usecase

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
	"fmt"
	"math/rand"
)

type paymentUsecase struct {
	repository repository.PaymentRepository
}

func NewPaymentUseCase(repo repository.PaymentRepository) *paymentUsecase {
	return &paymentUsecase{
		repository: repo,
	}
}

type PaymentUseCase interface {
	CreateSePayQR(ammount float64, user_id int) (models.CreateQR, error)
	SePayWebhook(transaction models.Transaction) error
}

func (p *paymentUsecase) CreateSePayQR(ammount float64, user_id int) (models.CreateQR, error) {
	// Generate description string
	description := fmt.Sprintf("AHV%07d", rand.Intn(10000000))

	// Generate the link string
	result := fmt.Sprintf("https://qr.sepay.vn/img?acc=%s&bank=%s&amount=%.2f&des=%s", "18896441", "ACB", ammount, description)

	// Construct the QR model
	qr := models.CreateQR{
		AccountNumber: "18896441",
		BankName:      "ACB",
		Amount:        ammount,
		Description:   description,
		Link:          result,
	}

	err := p.repository.CreateSePayQR(qr, user_id)
	if err != nil {
		return models.CreateQR{}, err
	}

	return qr, nil
}

func (p *paymentUsecase) SePayWebhook(transaction models.Transaction) error {

	err := p.repository.SePayWebhook(transaction)
	if err != nil {
		return err
	}

	return nil
}
