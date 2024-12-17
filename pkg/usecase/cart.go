package usecase

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
)

type CartUseCase interface {
	GetCart(user_id int) ([]models.CartItem, error)
	AddToCart(user_id, product_id int, quantity uint) (models.CartDetails, error)
	UpdateQuantityAdd(cart_id int, quantity uint) (models.CartDetails, error)
	UpdateQuantityLess(cart_id int, quantity uint) (models.CartDetails, error)
	UpdateQuantity(cart_id int, quantity uint) (models.CartDetails, error)
	RemoveFromCart(cart_id int) error

	// CheckOut(id int) (models.CheckOut, error)
}

type cartUseCase struct {
	repo              repository.CartRepository
	productRepository repository.ProductRepository
}

func NewCartUseCase(
	repo repository.CartRepository,
	productRepo repository.ProductRepository,
) *cartUseCase {
	return &cartUseCase{
		repo:              repo,
		productRepository: productRepo,
	}
}

func (i *cartUseCase) AddToCart(user_id, product_id int, quantity uint) (models.CartDetails, error) {

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

// func (i *cartUseCase) CheckOut(id int) (models.CheckOut, error) {

// 	// address, err := i.repo.GetAddresses(id)
// 	// if err != nil {
// 	// 	return models.CheckOut{}, err
// 	// }

// 	// payment, err := i.repo.GetPaymentOptions()
// 	// if err != nil {
// 	// 	return models.CheckOut{}, err
// 	// }

// 	// products, err := i.repo.GetCart(id)
// 	// if err != nil {
// 	// 	return models.CheckOut{}, err
// 	// }

// 	// var discountedPrice, totalPrice float64
// 	// for _, v := range products.Data {
// 	// 	discountedPrice += v.DiscountedPrice
// 	// 	totalPrice += v.Total
// 	// }

// 	var checkout models.CheckOut

// 	// checkout.CartID = products.ID
// 	// checkout.Addresses = address
// 	// checkout.Products = products.Data
// 	// checkout.PaymentMethods = payment
// 	// checkout.TotalPrice = totalPrice
// 	// checkout.DiscountedPrice = discountedPrice

// 	return checkout, nil
// }

func (i *cartUseCase) GetCart(user_id int) ([]models.CartItem, error) {

	cart, err := i.repo.GetCart(user_id)
	if err != nil {
		return []models.CartItem{}, err
	}

	return cart, nil
}

func (i *cartUseCase) RemoveFromCart(cart_id int) error {

	err := i.repo.RemoveFromCart(cart_id)
	if err != nil {
		return err
	}

	return nil
}

func (i *cartUseCase) UpdateQuantityAdd(cart_id int, quantity uint) (models.CartDetails, error) {

	result, err := i.repo.UpdateQuantityAdd(cart_id, quantity)
	if err != nil {
		return models.CartDetails{}, err
	}

	return result, nil
}

func (i *cartUseCase) UpdateQuantityLess(cart_id int, quantity uint) (models.CartDetails, error) {

	result, err := i.repo.UpdateQuantityLess(cart_id, quantity)
	if err != nil {
		return models.CartDetails{}, err
	}

	return result, nil
}

func (i *cartUseCase) UpdateQuantity(cart_id int, quantity uint) (models.CartDetails, error) {

	result, err := i.repo.UpdateQuantity(cart_id, quantity)
	if err != nil {
		return models.CartDetails{}, err
	}

	return result, nil
}

// func (u *cartUseCase) GetCart(id int) (models.GetCartResponse, error) {

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
// 	var product_names []string
// 	for i := range products {
// 		product_name, err := u.repo.FindProductNames(products[i])
// 		if err != nil {
// 			return models.GetCartResponse{}, errors.New(InternalError)
// 		}
// 		product_names = append(product_names, product_name)
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
// 	for i := range product_names {
// 		var get models.GetCart
// 		get.ID = products[i]
// 		get.ProductName = product_names[i]
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

// func (i *cartUseCase) RemoveFromCart(cart, product int) error {

// 	err := i.repo.RemoveFromCart(cart, product)
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }

// func (i *cartUseCase) UpdateQuantityAdd(id, inv int) error {

// 	err := i.repo.UpdateQuantityAdd(id, inv)
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }

// func (i *userUseCase) UpdateQuantityLess(id, inv int) error {

// 	err := i.repo.UpdateQuantityLess(id, inv)
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }
