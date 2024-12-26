package usecase

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
)

type OrderUseCase interface {
	PlaceOrder(order models.PlaceOrder) (models.Order, error)
}

type orderUseCase struct {
	repository repository.OrderRepository
}

func NewOrderUseCase(
	repo repository.OrderRepository,
) OrderUseCase {
	return &orderUseCase{
		repository: repo,
	}
}

func (or *orderUseCase) PlaceOrder(order models.PlaceOrder) (models.Order, error) {

	result, err := or.repository.PlaceOrder(order)
	if err != nil {
		return models.Order{}, err
	}

	for _, item := range order.CartCheckOut.CartItems {
		err := or.repository.PlaceOrderItem(item, result.ID)
		if err != nil {
			return models.Order{}, err
		}
	}

	return result, nil
}
