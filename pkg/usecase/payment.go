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
	CreateQR(ammount float64) (models.CreateQR, error)
}

func (p *paymentUsecase) CreateQR(ammount float64) (models.CreateQR, error) {
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

	return qr, nil
}
