package repository

import (
	"ahava/pkg/domain"
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	GetCart(user_id uint, cart_ids []uint) ([]models.CartItem, error)
	AddToCart(user_id uint, cart_item models.UpdateCartItem) (models.CartDetails, error)

	CheckIfItemIsAlreadyAdded(user_id, product_id uint, size string) (uint, error)
	UpdateQuantityAdd(user_id, cart_id, quantity uint) (models.CartDetails, error)
	UpdateQuantityLess(user_id, cart_id, quantity uint) (models.CartDetails, error)
	UpdateQuantity(user_id, cart_id, quantity uint) (models.CartDetails, error)
	RemoveFromCart(user_id, cart_id uint) error
}

type cartRepository struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{
		DB: db,
	}
}

func (r *cartRepository) GetCart(user_id uint, cart_ids []uint) ([]models.CartItem, error) {
	// Create a slice of cart items
	var cart []models.CartItem
	// Create a query to get the cart items
	query := r.DB.Model(&domain.CartItem{}).
		Joins("JOIN products p ON cart_items.product_id = p.id").
		Joins("JOIN prices pr ON pr.product_id = p.id AND pr.size = cart_items.size").
		Select(`cart_items.id, p.id as product_id, p.name, p.default_image, cart_items.quantity, cart_items.size, pr.original_price, pr.discount_price,
				(cart_items.quantity * pr.original_price) AS item_price, 
				(cart_items.quantity * pr.discount_price) AS item_discount_price`).
		Where("cart_items.user_id = ?", user_id)
	// If there are cart ids, add a where clause to the query
	if len(cart_ids) > 0 {
		query = query.Where("cart_items.id IN ?", cart_ids)
	}
	// Execute the query and scan the result into the cart slice
	if err := query.Find(&cart).Error; err != nil {
		return nil, err
	}
	// Return the cart slice
	return cart, nil
}

func (r *cartRepository) CheckIfItemIsAlreadyAdded(user_id, product_id uint, size string) (uint, error) {

	var cart_id uint

	err := r.DB.Model(&domain.CartItem{}).
		Select("id").
		Where("user_id=? AND product_id=? AND size=?", user_id, product_id, size).
		Scan(&cart_id).Error
	if err != nil {
		return 0, err
	}

	return cart_id, nil
}

func (r *cartRepository) RemoveFromCart(user_id, cart_id uint) error {
	// Delete the cart item
	result := r.DB.Where("id=? AND user_id=?", cart_id, user_id).Delete(&domain.CartItem{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrEntityNotFound
	}

	return nil
}

func (r *cartRepository) UpdateQuantityAdd(user_id, cart_id, quantity uint) (models.CartDetails, error) {

	var cartDetails models.CartDetails

	result := r.DB.Model(&domain.CartItem{}).
		Where("id=? AND user_id=?", cart_id, user_id).
		Update("quantity", gorm.Expr("quantity + ?", quantity)).
		Scan(&cartDetails)
	if result.Error != nil {
		return models.CartDetails{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.CartDetails{}, models.ErrEntityNotFound
	}

	return cartDetails, nil
}

func (r *cartRepository) UpdateQuantityLess(user_id, cart_id, quantity uint) (models.CartDetails, error) {

	var cartDetails models.CartDetails

	result := r.DB.
		Model(&domain.CartItem{}).
		Where("id=? AND user_id=?", cart_id, user_id).
		Update("quantity", gorm.Expr("quantity - ?", quantity)).
		Scan(&cartDetails)
	if result.Error != nil {
		return models.CartDetails{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.CartDetails{}, models.ErrEntityNotFound
	}

	return cartDetails, nil
}

func (r *cartRepository) UpdateQuantity(user_id, cart_id, quantity uint) (models.CartDetails, error) {

	var cartDetails models.CartDetails

	result := r.DB.
		Model(&domain.CartItem{}).
		Where("id=? AND user_id=?", cart_id, user_id).
		Update("quantity", quantity).
		Scan(&cartDetails)
	if result.Error != nil {
		return models.CartDetails{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.CartDetails{}, models.ErrEntityNotFound
	}

	return cartDetails, nil
}

func (r *cartRepository) AddToCart(user_id uint, i models.UpdateCartItem) (models.CartDetails, error) {

	cart_item := domain.CartItem{
		UserID:    user_id,
		ProductID: i.ProductID,
		Quantity:  i.Quantity,
		Size:      i.Size,
	}

	if err := r.DB.Create(&cart_item).Error; err != nil {
		return models.CartDetails{}, err
	}

	return models.CartDetails{
		ID:        cart_item.ID,
		UserID:    cart_item.UserID,
		ProductID: cart_item.ProductID,
		Quantity:  cart_item.Quantity,
		Size:      cart_item.Size,
	}, nil
}
