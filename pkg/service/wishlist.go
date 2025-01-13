package service

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
)

type WishlistService interface {
	AddToWishlist(user_id uint, product models.AddToWishlist) (models.Wishlist, error)
	RemoveFromWishlist(user_id, wishlist_id uint) error
	GetWishList(user_id uint, order_by string) ([]models.WishlistProduct, error)
}

type wishlistService struct {
	repository        repository.WishlistRepository
	productRepository repository.ProductRepository
}

func NewWishlistService(
	repo repository.WishlistRepository,
	productRepo repository.ProductRepository,
) WishlistService {
	return &wishlistService{
		repository:        repo,
		productRepository: productRepo,
	}
}

func (w *wishlistService) AddToWishlist(user_id uint, product models.AddToWishlist) (models.Wishlist, error) {

	exists, err := w.repository.CheckIfTheItemIsPresentAtWishlist(user_id, product.ProductID, product.Size)
	if err != nil {
		return models.Wishlist{}, err
	}

	if exists {
		result, err := w.repository.UpdateWishlist(user_id, product.ProductID, product.Size, false)
		if err != nil {
			return models.Wishlist{}, err
		}
		return result, nil
	} else {
		result, err := w.repository.AddToWishlist(user_id, product)
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

func (w *wishlistService) GetWishList(user_id uint, order_by string) ([]models.WishlistProduct, error) {

	products, err := w.repository.GetWishList(user_id, order_by)
	if err != nil {
		return []models.WishlistProduct{}, err
	}

	return products, nil
}
