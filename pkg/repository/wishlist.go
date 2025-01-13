package repository

import (
	"ahava/pkg/domain"
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type WishlistRepository interface {
	AddToWishlist(user_id uint, product models.AddToWishlist) (models.Wishlist, error)
	UpdateWishlist(user_id, product_id uint, size string, is_deleted bool) (models.Wishlist, error)
	UpdateRemoveFromWishlist(user_id, wishlist_id uint) error
	GetWishList(user_id uint, order_by string) ([]models.WishlistProduct, error)
	CheckIfTheItemIsPresentAtWishlist(user_id, product_id uint, size string) (bool, error)
	// CheckIfTheItemIsPresentAtCart(user_id, product_id uint) (bool, error)
}

type wishlistRepository struct {
	DB *gorm.DB
}

func NewWishlistRepository(DB *gorm.DB) WishlistRepository {
	return &wishlistRepository{DB}
}

func (r *wishlistRepository) AddToWishlist(user_id uint, product models.AddToWishlist) (models.Wishlist, error) {

	addWishlist := domain.Wishlist{
		UserID:    user_id,
		ProductID: product.ProductID,
		Size:      product.Size,
	}

	if err := r.DB.Create(&addWishlist).Error; err != nil {
		return models.Wishlist{}, err
	}

	return models.Wishlist{
		ID:        addWishlist.ID,
		UserID:    addWishlist.UserID,
		ProductID: addWishlist.ProductID,
		Size:      addWishlist.Size,
	}, nil
}

func (r *wishlistRepository) UpdateWishlist(user_id, product_id uint, size string, is_deleted bool) (models.Wishlist, error) {

	var updateWishlist models.Wishlist

	result := r.DB.
		Model(&domain.Wishlist{}).
		Where("user_id=? AND product_id=? AND size=?", user_id, product_id, size).
		Update("is_deleted", is_deleted).
		Scan(&updateWishlist)

	if result.Error != nil {
		return models.Wishlist{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.Wishlist{}, models.ErrEntityNotFound
	}

	return updateWishlist, nil
}

func (r *wishlistRepository) UpdateRemoveFromWishlist(user_id, wishlist_id uint) error {

	result := r.DB.
		Model(&domain.Wishlist{}).
		Where("id=? AND user_id=?", wishlist_id, user_id).
		Update("is_deleted", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrEntityNotFound
	}

	return nil
}

func (r *wishlistRepository) GetWishList(user_id uint, order_by string) ([]models.WishlistProduct, error) {

	var wishlistProducts []models.WishlistProduct

	query := r.DB.Model(&domain.Product{}).
		Select(`products.id AS product_id, products.name, products.default_image, wishlists.id,
				COUNT(wishlists.product_id) AS total_count, prices.original_price, prices.discount_price, wishlists.create_at`).
		Joins("JOIN wishlists ON wishlists.product_id = products.id").
		Joins("JOIN prices ON prices.product_id = products.id AND prices.size = wishlists.size"). // Join prices table
		Where("wishlists.is_deleted = false AND wishlists.user_id = ?", user_id).
		Group("wishlists.id, products.id, products.name, products.default_image, prices.discount_price, prices.original_price") // Group by necessary fields

	switch order_by {
	case "price_asc":
		query = query.Order("prices.discount_price ASC") // Ordering by the smallest discount price
	case "price_desc":
		query = query.Order("prices.discount_price DESC") // Ordering by the smallest discount price in descending order
	case "latest":
		query = query.Order("products.create_at DESC")
	case "most_favorite":
		query = query.Order("total_count DESC")
	case "most_viewed":
		query = query.Order("products.create_at DESC")
	default:
		query = query.Order("products.create_at DESC")
	}

	if err := query.Scan(&wishlistProducts).Error; err != nil {
		return nil, err
	}

	return wishlistProducts, nil
}

func (r *wishlistRepository) CheckIfTheItemIsPresentAtWishlist(user_id, product_id uint, size string) (bool, error) {

	var result int64

	if err := r.DB.Raw(`SELECT COUNT (*) FROM wishlists WHERE user_id=$1 AND product_id=$2 AND size=$3`,
		user_id, product_id, size).Scan(&result).Error; err != nil {
		return false, err
	}

	return result > 0, nil

}
