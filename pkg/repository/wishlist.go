package repository

import (
	"ahava/pkg/domain"
	errors "ahava/pkg/utils/errors"
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type WishlistRepository interface {
	AddToWishlist(user_id, product_id uint) (models.Wishlist, error)
	UpdateWishlist(user_id, product_id uint, is_deleted bool) (models.Wishlist, error)
	UpdateRemoveFromWishlist(user_id, wishlist_id uint) error
	GetWishList(user_id uint, order_by string) ([]models.WishlistProduct, error)
	CheckIfTheItemIsPresentAtWishlist(user_id, product_id uint) (bool, error)
	// CheckIfTheItemIsPresentAtCart(user_id, product_id uint) (bool, error)
}

type wishlistRepository struct {
	DB *gorm.DB
}

func NewWishlistRepository(DB *gorm.DB) WishlistRepository {
	return &wishlistRepository{DB}
}

func (r *wishlistRepository) AddToWishlist(user_id, product_id uint) (models.Wishlist, error) {

	var addWishlist models.Wishlist

	err := r.DB.Raw(`INSERT INTO wishlists (user_id,product_id) VALUES ($1,$2)`,
		user_id, product_id).Scan(&addWishlist).Error
	if err != nil {
		return models.Wishlist{}, err
	}

	return addWishlist, nil
}

func (r *wishlistRepository) UpdateWishlist(user_id, product_id uint, is_deleted bool) (models.Wishlist, error) {

	var updateWishlist models.Wishlist

	result := r.DB.
		Model(&domain.Wishlist{}).
		Where("user_id = ? AND product_id = ?", user_id, product_id).
		Update("is_deleted", is_deleted).
		Scan(&updateWishlist)

	if result.Error != nil {
		return models.Wishlist{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.Wishlist{}, errors.ErrEntityNotFound
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
		return errors.ErrEntityNotFound
	}

	return nil
}

func (r *wishlistRepository) GetWishList(user_id uint, order_by string) ([]models.WishlistProduct, error) {
	var wishlistProducts []models.WishlistProduct

	// Start building the query
	query := r.DB.Model(&domain.Product{}).
		Select(`
            products.id AS product_id,
            products.name,
            products.default_image,
            products.stock,
            wishlists.id,
            COUNT(wishlists.product_id) AS total_count`).
		Joins("JOIN wishlists ON wishlists.product_id = products.id").
		Where("wishlists.is_deleted = false").
		Group("products.id, wishlists.id, products.name, products.default_image, products.stock")

	// Add filtering for specific user
	query = query.Where("wishlists.user_id = ?", user_id)

	// Add dynamic ordering based on the order_by parameter
	switch order_by {
	case "price_asc":
		query = query.Order("products.price ASC")
	case "price_desc":
		query = query.Order("products.price DESC")
	case "latest":
		query = query.Order("wishlists.created_at DESC")
	case "most_favorite":
		query = query.Order("total_count DESC")
	case "most_viewed":
		query = query.Order("wishlists.created_at DESC")
	default:
		query = query.Order("wishlists.created_at DESC")
	}

	if err := query.Scan(&wishlistProducts).Error; err != nil {
		return nil, err
	}

	return wishlistProducts, nil
}

// func (r *wishlistRepository) CheckIfTheItemIsPresentAtCart(user_id, product_id uint) (bool, error) {

// 	var result int64

// 	if err := r.DB.Raw(`SELECT COUNT (*)
// 	 FROM line_items
// 	 JOIN carts ON carts.id = line_items.cart_id
// 	 JOIN users ON users.id = carts.user_id
// 	 WHERE users.id = $1
// 	 AND
// 	 line_items.product_id = $2`, user_id, product_id).Scan(&result).Error; err != nil {
// 		return false, err
// 	}

// 	return result > 0, nil

// }

func (r *wishlistRepository) CheckIfTheItemIsPresentAtWishlist(user_id, product_id uint) (bool, error) {

	var result int64

	if err := r.DB.Raw(`SELECT COUNT (*) FROM wishlists WHERE user_id=$1 AND product_id=$2`,
		user_id, product_id).Scan(&result).Error; err != nil {
		return false, err
	}

	return result > 0, nil

}
