package service

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
)

type CartService interface {
	GetCart(user_id uint, cart_ids []uint) ([]models.CartItem, error)
	AddToCart(user_id uint, cart_item models.UpdateCartItem) (models.CartDetails, error)
	UpdateQuantityAdd(user_id, cart_id uint, quantity uint) (models.CartDetails, error)
	UpdateQuantityLess(user_id, cart_id uint, quantity uint) (models.CartDetails, error)
	UpdateQuantity(user_id, cart_id uint, quantity uint) (models.CartDetails, error)
	RemoveFromCart(user_id, cart_id uint) error
	CheckOut(user_id uint, cart_ids []uint) (models.CheckOut, error)
}

type cartService struct {
	repo           repository.CartRepository
	userRepository repository.UserRepository
}

func NewCartService(
	repo repository.CartRepository,
	userRepository repository.UserRepository,
) CartService {
	return &cartService{
		repo:           repo,
		userRepository: userRepository,
	}
}

func (i *cartService) AddToCart(user_id uint, cart_item models.UpdateCartItem) (models.CartDetails, error) {

	cart_id, err := i.repo.CheckIfItemIsAlreadyAdded(user_id, cart_item.ProductID, cart_item.Size)
	if err != nil {
		return models.CartDetails{}, err
	}

	if cart_id != 0 {
		result, err := i.repo.UpdateQuantityAdd(user_id, cart_id, cart_item.Quantity)
		if err != nil {
			return models.CartDetails{}, err
		}
		return result, nil
	} else {
		result, err := i.repo.AddToCart(user_id, cart_item)
		if err != nil {
			return models.CartDetails{}, err
		}
		return result, nil
	}
}

func (i *cartService) CheckOut(user_id uint, cart_ids []uint) (models.CheckOut, error) {

	cartItems, err := i.repo.GetCart(user_id, cart_ids)
	if err != nil {
		return models.CheckOut{}, err
	}

	var discountedPrice, totalPrice uint64

	for _, v := range cartItems {
		totalPrice += v.ItemPrice
		discountedPrice += v.ItemDiscountPrice
	}

	var checkout models.CheckOut

	checkout.CartItems = cartItems
	checkout.TotalPrice = totalPrice
	checkout.TotalDiscountedPrice = discountedPrice

	return checkout, nil
}

func (i *cartService) GetCart(user_id uint, cart_ids []uint) ([]models.CartItem, error) {

	cartItems, err := i.repo.GetCart(user_id, cart_ids)
	if err != nil {
		return []models.CartItem{}, err
	}

	return cartItems, nil
}

func (i *cartService) RemoveFromCart(user_id, cart_id uint) error {

	err := i.repo.RemoveFromCart(user_id, cart_id)
	if err != nil {
		return err
	}

	return nil
}

func (i *cartService) UpdateQuantityAdd(user_id, cart_id uint, quantity uint) (models.CartDetails, error) {

	result, err := i.repo.UpdateQuantityAdd(user_id, cart_id, quantity)
	if err != nil {
		return models.CartDetails{}, err
	}

	return result, nil
}

func (i *cartService) UpdateQuantityLess(user_id, cart_id uint, quantity uint) (models.CartDetails, error) {

	result, err := i.repo.UpdateQuantityLess(user_id, cart_id, quantity)
	if err != nil {
		return models.CartDetails{}, err
	}

	if result.Quantity <= 0 {
		err := i.repo.RemoveFromCart(user_id, cart_id)
		if err != nil {
			return models.CartDetails{}, err
		}
	}

	return result, nil
}

func (i *cartService) UpdateQuantity(user_id, cart_id uint, quantity uint) (models.CartDetails, error) {

	result, err := i.repo.UpdateQuantity(user_id, cart_id, quantity)
	if err != nil {
		return models.CartDetails{}, err
	}

	if result.Quantity <= 0 {
		err := i.repo.RemoveFromCart(user_id, cart_id)
		if err != nil {
			return models.CartDetails{}, err
		}
	}

	return result, nil
}
