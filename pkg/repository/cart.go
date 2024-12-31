package repository

import (
	"ahava/pkg/domain"
	errors "ahava/pkg/utils/errors"
	"ahava/pkg/utils/models"

	"github.com/lib/pq"

	"gorm.io/gorm"
)

type CartRepository interface {
	GetCart(user_id uint, cart_ids []uint) ([]models.CartItem, error)
	AddToCart(user_id, product_id uint, quantity uint) (models.CartDetails, error)
	CheckIfItemIsAlreadyAdded(user_id, product_id uint) (uint, error)
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

// func (r *cartRepository) GetAddresses(id uint) ([]models.Address, error) {

// 	var addresses []models.Address

// 	if err := r.DB.Raw("SELECT * FROM addresses WHERE user_id=$1", id).Scan(&addresses).Error; err != nil {
// 		return []models.Address{}, err
// 	}

// 	return addresses, nil

// }

func (r *cartRepository) GetCart(user_id uint, cart_ids []uint) ([]models.CartItem, error) {
	var cart []models.CartItem

	query := `
        SELECT ci.id AS cart_id, p.id as product_id, p.name, p.default_image, p.price, ci.quantity, 
               (ci.quantity * p.price) AS item_price 
        FROM cart_items ci JOIN products p ON ci.product_id = p.id WHERE ci.user_id = $1`

	if len(cart_ids) > 0 {
		query += " AND ci.id = ANY($2)"
	}

	var err error
	if len(cart_ids) > 0 {
		err = r.DB.Raw(query, user_id, pq.Array(cart_ids)).Scan(&cart).Error
	} else {
		err = r.DB.Raw(query, user_id).Scan(&cart).Error
	}

	if err != nil {
		return []models.CartItem{}, err
	}

	return cart, nil
}

// func (r *cartRepository) GetCartId(user_id uint) (int, error) {

// 	var id uint

// 	if err := r.DB.Raw("SELECT id FROM carts WHERE user_id=?", user_id).Scan(&id).Error; err != nil {
// 		return 0, err
// 	}

// 	return id, nil

// }

// func (i *cartRepository) CreateNewCart(user_id uint) (int, error) {
// 	var id uint
// 	err := i.DB.Exec(`
// 		INSERT INTO carts (user_id)
// 		VALUES ($1)`, user_id).Error
// 	if err != nil {
// 		return 0, err
// 	}

// 	if err := i.DB.Raw("select id from carts where user_id=?", user_id).Scan(&id).Error; err != nil {
// 		return 0, err
// 	}

// 	return id, nil
// }

// func (i *cartRepository) AddLineItems(cart_id, product_id uint) error {

// 	err := i.DB.Exec(`
// 		INSERT INTO line_items (cart_id,product_id)
// 		VALUES ($1,$2)`, cart_id, product_id).Error
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (r *cartRepository) CheckIfItemIsAlreadyAdded(user_id, product_id uint) (uint, error) {

	var count uint

	if err := r.DB.Raw("SELECT id FROM cart_items WHERE user_id = $1 AND product_id=$2",
		user_id, product_id).Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *cartRepository) RemoveFromCart(user_id, cart_id uint) error {

	err := r.DB.Exec(`DELETE FROM cart_items WHERE id=$1 AND user_id=$2`,
		cart_id, user_id).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *cartRepository) UpdateQuantityAdd(user_id, cart_id, quantity uint) (models.CartDetails, error) {

	var cartDetails models.CartDetails

	result := r.DB.
		Model(&domain.CartItems{}).
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
		Model(&domain.CartItems{}).
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
		Model(&domain.CartItems{}).
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

func (r *cartRepository) AddToCart(user_id, product_id, quantity uint) (models.CartDetails, error) {

	var cartDetails models.CartDetails

	err := r.DB.Raw(`INSERT INTO cart_items (user_id,product_id,quantity) VALUES ($1,$2,$3) RETURNING id, user_id, product_id, quantity`,
		user_id, product_id, quantity).Scan(&cartDetails).Error
	if err != nil {
		return models.CartDetails{}, err
	}

	return cartDetails, nil
}
