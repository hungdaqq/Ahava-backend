package usecase

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
	"fmt"
	"math/rand"
)

type paymentUsecase struct {
	repository      repository.PaymentRepository
	orderRepository repository.OrderRepository
}

func NewPaymentUseCase(repo repository.PaymentRepository, orderRepository repository.OrderRepository) *paymentUsecase {
	return &paymentUsecase{
		repository:      repo,
		orderRepository: orderRepository,
	}
}

type PaymentUseCase interface {
	CreateSePayQR(user_id, order_id uint, amount uint64) (models.CreateQR, error)
	Webhook(transaction models.Transaction) error
}

func (p *paymentUsecase) CreateSePayQR(user_id, order_id uint, amount uint64) (models.CreateQR, error) {

	description := fmt.Sprintf("AHV%07d", rand.Intn(10000000))

	result := fmt.Sprintf("https://qr.sepay.vn/img?acc=%s&bank=%s&amount=%d&des=%s&template=compact", "18896441", "ACB", amount, description)

	qr := models.CreateQR{
		OrderID:       order_id,
		AccountNumber: "18896441",
		BankName:      "ACB",
		Amount:        amount,
		Description:   description,
		Link:          result,
	}

	err := p.repository.CreateQR(qr, user_id)
	if err != nil {
		return models.CreateQR{}, err
	}

	return qr, nil
}

func (p *paymentUsecase) Webhook(transaction models.Transaction) error {

	// Save the transaction
	transaction, err := p.repository.SaveTransaction(transaction)
	if err != nil {
		return err
	}

	// Get the order details
	order, err := p.orderRepository.GetOrderDetails(transaction.OrderID)
	if err != nil {
		return err
	}

	// Check if the order has been paid: if the final price is greater than or equal to the transfer amount
	if order.FinalPrice >= transaction.TransferAmount {
		_, err = p.orderRepository.UpdateOrder(transaction.OrderID, models.Order{PaymentStatus: "PAID"})
		if err != nil {
			return err
		}
	} else {
		_, err = p.orderRepository.UpdateOrder(transaction.OrderID, models.Order{PaymentStatus: "INCOMPLETE"})
		if err != nil {
			return err
		}
	}

	return nil
}
