package usecase

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
)

type OrderUseCase interface {
	PlaceOrder(order models.PlaceOrder) (models.Order, error)
	GetOrderDetails(user_id, order_id uint) (models.Order, error)
	UpdateOrder(order_id uint, updateOrder models.Order) (models.Order, error)
}

type orderUseCase struct {
	repository  repository.OrderRepository
	cartUseCase CartUseCase
}

func NewOrderUseCase(
	repo repository.OrderRepository,
	cartUseCase CartUseCase,
) *orderUseCase {
	return &orderUseCase{
		repository:  repo,
		cartUseCase: cartUseCase,
	}
}

func (or *orderUseCase) PlaceOrder(placeOrder models.PlaceOrder) (models.Order, error) {

	checkout, err := or.cartUseCase.CheckOut(placeOrder.UserID, placeOrder.CartIDs)
	if err != nil {
		return models.Order{}, err
	}

	order, err := or.repository.PlaceOrder(placeOrder, checkout.TotalDiscountedPrice)
	if err != nil {
		return models.Order{}, err
	}

	for _, item := range checkout.CartItems {

		err := or.repository.PlaceOrderItem(order.ID, item)
		if err != nil {
			return models.Order{}, err
		}
	}

	return order, nil
}

func (or *orderUseCase) GetOrderDetails(user_id, order_id uint) (models.Order, error) {

	result, err := or.repository.GetOrderDetails(user_id, order_id)
	if err != nil {
		return models.Order{}, err
	}

	return result, nil
}

func (or *orderUseCase) UpdateOrder(order_id uint, updateOrder models.Order) (models.Order, error) {

	result, err := or.repository.UpdateOrder(order_id, updateOrder)
	if err != nil {
		return models.Order{}, err
	}

	return result, nil
}

// func (or *orderUseCase) UpdateOrderPaidStatus(order_id uint, status string) error {

// 	err := or.repository.UpdateOrderPaidStatus(order_id, status)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
