package usecase

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
)

type CartUseCase interface {
	GetCart(user_id uint, cart_ids []uint) ([]models.CartItem, error)
	AddToCart(user_id, product_id uint, quantity uint) (models.CartDetails, error)
	UpdateQuantityAdd(cart_id uint, quantity uint) (models.CartDetails, error)
	UpdateQuantityLess(cart_id uint, quantity uint) (models.CartDetails, error)
	UpdateQuantity(cart_id uint, quantity uint) (models.CartDetails, error)
	RemoveFromCart(cart_id uint) error

	CheckOut(user_id uint, cart_ids []uint) (models.CheckOut, error)
}

type cartUseCase struct {
	repo            repository.CartRepository
	userRepository  repository.UserRepository
	offerRepository repository.OfferRepository
}

func NewCartUseCase(
	repo repository.CartRepository,
	userRepository repository.UserRepository,
	offerRepository repository.OfferRepository,
) *cartUseCase {
	return &cartUseCase{
		repo:            repo,
		userRepository:  userRepository,
		offerRepository: offerRepository,
	}
}

func (i *cartUseCase) AddToCart(user_id, product_id uint, quantity uint) (models.CartDetails, error) {

	cart_id, err := i.repo.CheckIfItemIsAlreadyAdded(user_id, product_id)
	if err != nil {
		return models.CartDetails{}, err
	}

	if cart_id != 0 {
		result, err := i.repo.UpdateQuantityAdd(cart_id, quantity)
		if err != nil {
			return models.CartDetails{}, err
		}
		return result, nil
	} else {
		result, err := i.repo.AddToCart(user_id, product_id, quantity)
		if err != nil {
			return models.CartDetails{}, err
		}
		return result, nil
	}
}

func (i *cartUseCase) CheckOut(user_id uint, cart_ids []uint) (models.CheckOut, error) {

	cartItems, err := i.repo.GetCart(user_id, cart_ids)

	var discountedPrice, totalPrice uint64

	for idx := range cartItems {
		offerPercentage, err := i.offerRepository.FindOfferRate(cartItems[idx].ProductID)
		if err != nil {
			return models.CheckOut{}, err
		}
		if offerPercentage > 0 {
			cartItems[idx].ItemDiscountedPrice = cartItems[idx].ItemPrice - (cartItems[idx].ItemPrice*uint64(offerPercentage))/100
		} else {
			cartItems[idx].ItemDiscountedPrice = cartItems[idx].ItemPrice
		}

		totalPrice += cartItems[idx].ItemPrice
		discountedPrice += cartItems[idx].ItemDiscountedPrice
	}

	if err != nil {
		return models.CheckOut{}, err
	}

	var checkout models.CheckOut

	checkout.CartItems = cartItems
	checkout.TotalPrice = totalPrice
	checkout.TotalDiscountedPrice = discountedPrice

	return checkout, nil
}

func (i *cartUseCase) GetCart(user_id uint, cart_ids []uint) ([]models.CartItem, error) {

	cartItems, err := i.repo.GetCart(user_id, cart_ids)
	if err != nil {
		return []models.CartItem{}, err
	}

	for idx := range cartItems {
		offerPercentage, err := i.offerRepository.FindOfferRate(cartItems[idx].ProductID)
		if err != nil {
			return []models.CartItem{}, err
		}
		if offerPercentage > 0 {
			cartItems[idx].ItemDiscountedPrice = cartItems[idx].ItemPrice - (cartItems[idx].ItemPrice*uint64(offerPercentage))/100
		} else {
			cartItems[idx].ItemDiscountedPrice = cartItems[idx].ItemPrice
		}
	}

	return cartItems, nil
}

func (i *cartUseCase) RemoveFromCart(cart_id uint) error {

	err := i.repo.RemoveFromCart(cart_id)
	if err != nil {
		return err
	}

	return nil
}

func (i *cartUseCase) UpdateQuantityAdd(cart_id uint, quantity uint) (models.CartDetails, error) {

	result, err := i.repo.UpdateQuantityAdd(cart_id, quantity)
	if err != nil {
		return models.CartDetails{}, err
	}

	return result, nil
}

func (i *cartUseCase) UpdateQuantityLess(cart_id uint, quantity uint) (models.CartDetails, error) {

	result, err := i.repo.UpdateQuantityLess(cart_id, quantity)
	if err != nil {
		return models.CartDetails{}, err
	}

	if result.Quantity <= 0 {
		err := i.repo.RemoveFromCart(cart_id)
		if err != nil {
			return models.CartDetails{}, err
		}
	}

	return result, nil
}

func (i *cartUseCase) UpdateQuantity(cart_id uint, quantity uint) (models.CartDetails, error) {

	result, err := i.repo.UpdateQuantity(cart_id, quantity)
	if err != nil {
		return models.CartDetails{}, err
	}

	if result.Quantity <= 0 {
		err := i.repo.RemoveFromCart(cart_id)
		if err != nil {
			return models.CartDetails{}, err
		}
	}

	return result, nil
}

// func (u *cartUseCase) GetCart(id uint) (models.GetCartResponse, error) {

// 	//find cart id
// 	cart_id, err := u.repo.GetCartID(id)
// 	if err != nil {
// 		return models.GetCartResponse{}, errors.New(InternalError)
// 	}
// 	//find products inide cart
// 	products, err := u.repo.GetProductsInCart(cart_id)
// 	if err != nil {
// 		return models.GetCartResponse{}, errors.New(InternalError)
// 	}
// 	//find product names
// 	var names []string
// 	for i := range products {
// 		name, err := u.repo.FindProductNames(products[i])
// 		if err != nil {
// 			return models.GetCartResponse{}, errors.New(InternalError)
// 		}
// 		names = append(names, name)
// 	}

// 	//find quantity
// 	var quantity []int
// 	for i := range products {
// 		q, err := u.repo.FindCartQuantity(cart_id, products[i])
// 		if err != nil {
// 			return models.GetCartResponse{}, errors.New(InternalError)
// 		}
// 		quantity = append(quantity, q)
// 	}

// 	var price []float64
// 	for i := range products {
// 		q, err := u.repo.FindPrice(products[i])
// 		if err != nil {
// 			return models.GetCartResponse{}, errors.New(InternalError)
// 		}
// 		price = append(price, q)
// 	}

// 	var images []string
// 	var stocks []int

// 	for _, v := range products {
// 		image, err := u.repo.FindProductImage(v)
// 		if err != nil {
// 			return models.GetCartResponse{}, errors.New(InternalError)
// 		}

// 		stock, err := u.repo.FindStock(v)
// 		if err != nil {
// 			return models.GetCartResponse{}, errors.New(InternalError)
// 		}

// 		images = append(images, image)
// 		stocks = append(stocks, stock)
// 	}

// 	var categories []int
// 	for i := range products {
// 		c, err := u.repo.FindCategory(products[i])
// 		if err != nil {
// 			return models.GetCartResponse{}, errors.New(InternalError)
// 		}
// 		categories = append(categories, c)
// 	}

// 	var getcart []models.GetCart
// 	for i := range names {
// 		var get models.GetCart
// 		get.ID = products[i]
// 		get.ProductName = names[i]
// 		get.Image = images[i]
// 		get.Category_id = categories[i]
// 		get.Quantity = quantity[i]
// 		get.Total = price[i]
// 		get.StockAvailable = stocks[i]
// 		get.DiscountedPrice = 0

// 		getcart = append(getcart, get)
// 	}

// 	//find offers
// 	var offers []int
// 	for i := range categories {
// 		c, err := u.repo.FindofferPercentage(categories[i])
// 		if err != nil {
// 			return models.GetCartResponse{}, errors.New(InternalError)
// 		}
// 		offers = append(offers, c)
// 	}

// 	//find discounted price
// 	for i := range getcart {
// 		getcart[i].DiscountedPrice = (getcart[i].Total) - (getcart[i].Total * float64(offers[i]) / 100)
// 	}

// 	var response models.GetCartResponse
// 	response.ID = cart_id
// 	response.Data = getcart

// 	//then return in appropriate format

// 	return response, nil

// }
