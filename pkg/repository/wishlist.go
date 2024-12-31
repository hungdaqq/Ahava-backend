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
	GetWishList(user_id uint) ([]models.Products, error)
	CheckIfTheItemIsPresentAtWishlist(user_id, product_id uint) (bool, error)
	// CheckIfTheItemIsPresentAtCart(user_id, product_id uint) (bool, error)
}

type wishlistRepository struct {
	DB *gorm.DB
}

func NewWishlistRepository(db *gorm.DB) WishlistRepository {
	return &wishlistRepository{
		DB: db,
	}
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

func (r *wishlistRepository) GetWishList(id uint) ([]models.Products, error) {
	var productDetails []models.Products

	query := `
        SELECT products.id,
               products.category_id,
               products.name,
               products.default_image,
               products.size,
               products.stock,
               products.price
        FROM products
        JOIN wishlists ON wishlists.product_id = products.id
        WHERE wishlists.user_id = ? AND wishlists.is_deleted = false
    `

	if err := r.DB.Raw(query, id).Scan(&productDetails).Error; err != nil {
		// Log or handle the error appropriately.
		return nil, err
	}

	return productDetails, nil
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
