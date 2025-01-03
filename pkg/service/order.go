package service

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
)

type OrderService interface {
	PlaceOrder(order models.PlaceOrder) (models.Order, error)
	GetOrderDetails(user_id, order_id uint) (models.Order, error)
	UpdateOrder(order_id uint, updateOrder models.Order) (models.Order, error)
}

type orderService struct {
	repository  repository.OrderRepository
	cartService CartService
}

func NewOrderService(repo repository.OrderRepository, cartService CartService) OrderService {
	return &orderService{
		repository:  repo,
		cartService: cartService,
	}
}

func (or *orderService) PlaceOrder(placeOrder models.PlaceOrder) (models.Order, error) {

	checkout, err := or.cartService.CheckOut(placeOrder.UserID, placeOrder.CartIDs)
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

func (or *orderService) GetOrderDetails(user_id, order_id uint) (models.Order, error) {

	result, err := or.repository.GetOrderDetails(user_id, order_id)
	if err != nil {
		return models.Order{}, err
	}

	return result, nil
}

func (or *orderService) UpdateOrder(order_id uint, updateOrder models.Order) (models.Order, error) {

	result, err := or.repository.UpdateOrder(order_id, updateOrder)
	if err != nil {
		return models.Order{}, err
	}

	return result, nil
}

// func (or *orderService) UpdateOrderPaidStatus(order_id uint, status string) error {

// 	err := or.repository.UpdateOrderPaidStatus(order_id, status)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
