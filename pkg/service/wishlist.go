package service

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
)

type WishlistService interface {
	AddToWishlist(user_id, product_id uint) (models.Wishlist, error)
	RemoveFromWishlist(user_id, wishlist_id uint) error
	GetWishList(user_id uint) ([]models.WishlistProduct, error)
}

type wishlistService struct {
	repository      repository.WishlistRepository
	offerRepository repository.OfferRepository
	// offerRepo  repository.OfferRepository
}

func NewWishlistService(
	repo repository.WishlistRepository,
	offer repository.OfferRepository,
) WishlistService {
	return &wishlistService{
		repository:      repo,
		offerRepository: offer,
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

func (w *wishlistService) RemoveFromWishlist(user_id, wishlist_id uint) error {

	if err := w.repository.UpdateRemoveFromWishlist(user_id, wishlist_id); err != nil {
		return err
	}

	return nil
}

func (w *wishlistService) GetWishList(user_id uint) ([]models.WishlistProduct, error) {

	products, err := w.repository.GetWishList(user_id)
	if err != nil {
		return []models.WishlistProduct{}, err
	}

	for idx := range products {
		offerPercentage, err := w.offerRepository.FindOfferRate(products[idx].ProductID)
		if err != nil {
			return nil, err
		}
		if offerPercentage > 0 {
			products[idx].DiscountedPrice = products[idx].Price - (products[idx].Price*uint64(offerPercentage))/100
		} else {
			products[idx].DiscountedPrice = products[idx].Price
		}
	}

	// //loop inside products and then calculate discounted price of each then return
	// for j := range productDetails {
	// 	discount_percentage, err := w.offerRepo.FindDiscountPercentage(productDetails[j].CategoryID)
	// 	if err != nil {
	// 		return []models.Product{}, errors.New("there was some error in finding the discounted prices")
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

	return products, nil
}
