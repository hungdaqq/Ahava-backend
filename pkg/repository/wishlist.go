package repository

import (
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type WishlistRepository interface {
	AddToWishlist(user_id, product_id uint) (models.Wishlist, error)
	RemoveFromWishlist(user_id, product_id uint) error
	GetWishList(user_id uint) ([]models.Products, error)

	CheckIfTheItemIsPresentAtWishlist(user_id, product_id uint) (bool, error)
	// CheckIfTheItemIsPresentAtCart(user_id, product_id uint) (bool, error)
}

type wishlistRepository struct {
	DB *gorm.DB
}

func NewWishlistRepository(db *gorm.DB) *wishlistRepository {
	return &wishlistRepository{
		DB: db,
	}
}

func (w *wishlistRepository) AddToWishlist(user_id, product_id uint) (models.Wishlist, error) {

	var addWishlist models.Wishlist

	err := w.DB.Exec(`INSERT INTO wishlists (user_id,product_id) VALUES ($1,$2)`,
		user_id, product_id).Scan(&addWishlist).Error
	if err != nil {
		return models.Wishlist{}, err
	}

	return addWishlist, nil
}

func (w *wishlistRepository) RemoveFromWishlist(user_id, product_id uint) error {

	err := w.DB.Exec("DELETE FROM wishlists WHERE product_id=$1 AND user_id=$2",
		user_id, product_id).Error
	if err != nil {
		return err
	}

	return nil
}

func (w *wishlistRepository) GetWishList(id uint) ([]models.Products, error) {
	var productDetails []models.Products

	query := `
        SELECT products.id,
               products.category_id,
               products.product_name,
               products.image,
               products.size,
               products.stock,
               products.price
        FROM products
        JOIN wishlists ON wishlists.product_id = products.id
        WHERE wishlists.user_id = ? AND wishlists.is_deleted = false
    `

	if err := w.DB.Raw(query, id).Scan(&productDetails).Error; err != nil {
		// Log or handle the error appropriately.
		return nil, err
	}

	return productDetails, nil
}

// func (w *wishlistRepository) CheckIfTheItemIsPresentAtCart(user_id, product_id uint) (bool, error) {

// 	var result int64

// 	if err := w.DB.Raw(`SELECT COUNT (*)
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

func (w *wishlistRepository) CheckIfTheItemIsPresentAtWishlist(user_id, product_id uint) (bool, error) {

	var result int64

	if err := w.DB.Raw(`SELECT COUNT (*) FROM wishlists WHERE user_id=$1 AND product_id=$2`,
		user_id, product_id).Scan(&result).Error; err != nil {
		return false, err
	}

	return result > 0, nil

}
