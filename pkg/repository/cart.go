package repository

import (
	"ahava/pkg/domain"
	errors "ahava/pkg/utils/errors"
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

	// GetAddresses(id uint) ([]models.Address, error)
	// GetPaymentOptions() ([]models.PaymentMethod, error)
	// GetCartId(user_id uint) (int, error)
	// CreateNewCart(user_id uint) (int, error)

	// AddLineItems(cart_id, product_id uint) error
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

	var cart []models.CartItem

	query := r.DB.Model(&models.CartItem{}).
		Joins("JOIN products p ON cart_items.product_id = p.id").
		Select("cart_items.id AS cart_id, p.id as product_id, p.name, p.default_image, cart_items.quantity, (cart_items.quantity * p.price) AS item_price").
		Where("cart_items.user_id = ?", user_id)

	if len(cart_ids) > 0 {
		query = query.Where("cart_items.id IN ?", cart_ids)
	}

	// Execute the query and scan the result into the cart slice
	if err := query.Find(&cart).Error; err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *cartRepository) CheckIfItemIsAlreadyAdded(user_id, product_id uint, size string) (uint, error) {

	var cart_id uint

	err := r.DB.Model(&domain.CartItem{}).
		Select("id").
		Where("user_id = ? AND product_id = ? AND size = ?", user_id, product_id, size).
		Scan(&cart_id).Error

	if err != nil {
		return 0, err
	}

	return cart_id, nil
}

func (r *cartRepository) RemoveFromCart(user_id, cart_id uint) error {

	result := r.DB.Exec(`DELETE FROM cart_items WHERE id=$1 AND user_id=$2`,
		cart_id, user_id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrEntityNotFound
	}

	return nil
}

func (r *cartRepository) UpdateQuantityAdd(user_id, cart_id, quantity uint) (models.CartDetails, error) {

	var cartDetails models.CartDetails

	result := r.DB.
		Model(&domain.CartItem{}).
		Where("id=? AND user_id=?", cart_id, user_id).
		Update("quantity", gorm.Expr("quantity + ?", quantity)).
		Scan(&cartDetails)

	if result.Error != nil {
		return models.CartDetails{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.CartDetails{}, errors.ErrEntityNotFound
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
		return models.CartDetails{}, errors.ErrEntityNotFound
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
		return models.CartDetails{}, errors.ErrEntityNotFound
	}

	return cartDetails, nil
}

func (r *cartRepository) AddToCart(user_id uint, cart_item models.UpdateCartItem) (models.CartDetails, error) {

	newCartItem := domain.CartItem{
		UserID:    user_id,
		ProductID: cart_item.ProductID,
		Quantity:  cart_item.Quantity,
		Size:      cart_item.Size,
	}

	if err := r.DB.Create(&newCartItem).Error; err != nil {
		return models.CartDetails{}, err
	}

	return models.CartDetails{
		ID:        newCartItem.ID,
		UserID:    newCartItem.UserID,
		ProductID: newCartItem.ProductID,
		Quantity:  newCartItem.Quantity,
	}, nil
}
