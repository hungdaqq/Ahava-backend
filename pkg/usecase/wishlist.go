package usecase

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
	"errors"
)

type WishlistUseCase interface {
	AddToWishlist(user_id, product_id int) (models.Wishlist, error)
	RemoveFromWishlist(user_id, product_id int) error
	GetWishList(user_id int) ([]models.Products, error)
}

type wishlistUseCase struct {
	repository repository.WishlistRepository
	// offerRepo  repository.OfferRepository
}

func NewWishlistUseCase(
	repo repository.WishlistRepository,
	// offer repository.OfferRepository,
) *wishlistUseCase {
	return &wishlistUseCase{
		repository: repo,
		// offerRepo:  offer,
	}
}

func (w *wishlistUseCase) AddToWishlist(user_id, product_id int) (models.Wishlist, error) {

	exists, err := w.repository.CheckIfTheItemIsPresentAtWishlist(user_id, product_id)
	if err != nil {
		return models.Wishlist{}, err
	}

	if exists {
		return models.Wishlist{}, errors.New("item already exists in wishlist")
	}

	result, err := w.repository.AddToWishlist(user_id, product_id)
	if err != nil {
		return models.Wishlist{}, errors.New("could not add to wishlist")
	}

	return result, nil
}

func (w *wishlistUseCase) RemoveFromWishlist(user_id, product_id int) error {

	if err := w.repository.RemoveFromWishlist(user_id, product_id); err != nil {
		return errors.New("could not remove from wishlist")
	}

	return nil
}

func (w *wishlistUseCase) GetWishList(id int) ([]models.Products, error) {

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
