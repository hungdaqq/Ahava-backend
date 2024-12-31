package service

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
	"errors"
)

type WishlistService interface {
	AddToWishlist(user_id, product_id uint) (models.Wishlist, error)
	RemoveFromWishlist(user_id, product_id uint) error
	GetWishList(user_id uint) ([]models.Products, error)
}

type wishlistService struct {
	repository repository.WishlistRepository
	// offerRepo  repository.OfferRepository
}

func NewWishlistService(
	repo repository.WishlistRepository,
	// offer repository.OfferRepository,
) WishlistService {
	return &wishlistService{
		repository: repo,
		// offerRepo:  offer,
	}
}

func (w *wishlistService) AddToWishlist(user_id, product_id uint) (models.Wishlist, error) {

	exists, err := w.repository.CheckIfTheItemIsPresentAtWishlist(user_id, product_id)
	if err != nil {
		return models.Wishlist{}, err
	}

	if exists {
		result, err := w.repository.UpdateWishlist(user_id, product_id, false)
		if err != nil {
			return models.Wishlist{}, err
		}
		return result, nil
	} else {
		result, err := w.repository.AddToWishlist(user_id, product_id)
		if err != nil {
			return models.Wishlist{}, err
		}
		return result, nil
	}
}

func (w *wishlistService) RemoveFromWishlist(user_id, product_id uint) error {

	if _, err := w.repository.UpdateWishlist(user_id, product_id, true); err != nil {
		return errors.New("could not remove from wishlist")
	}

	return nil
}

func (w *wishlistService) GetWishList(id uint) ([]models.Products, error) {

	productDetails, err := w.repository.GetWishList(id)
	if err != nil {
		return []models.Products{}, err
	}

	// //loop inside products and then calculate discounted price of each then return
	// for j := range productDetails {
	// 	discount_percentage, err := w.offerRepo.FindDiscountPercentage(productDetails[j].CategoryID)
	// 	if err != nil {
	// 		return []models.Products{}, errors.New("there was some error in finding the discounted prices")
	// 	}
	// 	var discount float64

	// 	if discount_percentage > 0 {
	// 		discount = (productDetails[j].Price * float64(discount_percentage)) / 100
	// 	}

	// 	productDetails[j].DiscountedPrice = productDetails[j].Price - discount

	// 	productDetails[j].IfPresentAtWishlist, err = w.repository.CheckIfTheItemIsPresentAtWishlist(id, int(productDetails[j].ID))
	// 	if err != nil {
	// 		return []models.Product{}, errors.New("error while checking ")
	// 	}

	// 	productDetails[j].IfPresentAtCart, err = w.repository.CheckIfTheItemIsPresentAtCart(id, int(productDetails[j].ID))
	// 	if err != nil {
	// 		return []models.ProductDetails{}, errors.New("error while checking ")
	// 	}

	// }

	return productDetails, nil
}
