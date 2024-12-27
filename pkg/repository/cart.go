package repository

import (
	"ahava/pkg/domain"
	"ahava/pkg/utils/models"
	"errors"

	"github.com/lib/pq"

	"gorm.io/gorm"
)

type CartRepository interface {
	GetCart(user_id uint, cart_ids []uint) ([]models.CartItem, error)
	AddToCart(user_id, product_id uint, quantity uint) (models.CartDetails, error)
	CheckIfItemIsAlreadyAdded(user_id, product_id uint) (uint, error)
	UpdateQuantityAdd(cart_id uint, quantity uint) (models.CartDetails, error)
	UpdateQuantityLess(cart_id uint, quantity uint) (models.CartDetails, error)
	UpdateQuantity(cart_id uint, quantity uint) (models.CartDetails, error)

	RemoveFromCart(cart_id uint) error

	// GetAddresses(id uint) ([]models.Address, error)
	// GetPaymentOptions() ([]models.PaymentMethod, error)
	// GetCartId(user_id uint) (int, error)
	// CreateNewCart(user_id uint) (int, error)

	// AddLineItems(cart_id, product_id uint) error
}

type cartRepository struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) *cartRepository {
	return &cartRepository{
		DB: db,
	}
}

// func (ad *cartRepository) GetAddresses(id uint) ([]models.Address, error) {

// 	var addresses []models.Address

// 	if err := ad.DB.Raw("SELECT * FROM addresses WHERE user_id=$1", id).Scan(&addresses).Error; err != nil {
// 		return []models.Address{}, err
// 	}

// 	return addresses, nil

// }

func (ad *cartRepository) GetCart(user_id uint, cart_ids []uint) ([]models.CartItem, error) {
	var cart []models.CartItem

	query := `
        SELECT ci.id AS cart_id, p.id as product_id, p.product_name, p.image, p.price, ci.quantity, 
               (ci.quantity * p.price) AS item_price 
        FROM cart_items ci JOIN products p ON ci.product_id = p.id WHERE ci.user_id = $1`

	if len(cart_ids) > 0 {
		query += " AND ci.id = ANY($2)"
	}

	var err error
	if len(cart_ids) > 0 {
		err = ad.DB.Raw(query, user_id, pq.Array(cart_ids)).Scan(&cart).Error
	} else {
		err = ad.DB.Raw(query, user_id).Scan(&cart).Error
	}

	if err != nil {
		return []models.CartItem{}, err
	}

	return cart, nil
}

// func (ad *cartRepository) GetCartId(user_id uint) (int, error) {

// 	var id uint

// 	if err := ad.DB.Raw("SELECT id FROM carts WHERE user_id=?", user_id).Scan(&id).Error; err != nil {
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

func (ad *cartRepository) CheckIfItemIsAlreadyAdded(user_id, product_id uint) (uint, error) {

	var count uint

	if err := ad.DB.Raw("SELECT id FROM cart_items WHERE user_id = $1 AND product_id=$2",
		user_id, product_id).Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (ad *cartRepository) RemoveFromCart(cart_id uint) error {

	err := ad.DB.Exec(`DELETE FROM cart_items WHERE id=$1`,
		cart_id).Error
	if err != nil {
		return err
	}

	return nil

}

func (ad *cartRepository) UpdateQuantityAdd(cart_id uint, quantity uint) (models.CartDetails, error) {

	var cartItem domain.CartItems

	result := ad.DB.Model(&cartItem).Where("id = ?", cart_id).Update("quantity", gorm.Expr("quantity + ?", quantity))

	if result.Error != nil {
		return models.CartDetails{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.CartDetails{}, errors.New("cart item not found")
	}

	return models.CartDetails{
		ID:        cartItem.ID,
		UserID:    cartItem.UserID,
		ProductID: cartItem.ProductID,
		Quantity:  cartItem.Quantity,
	}, nil
}

func (ad *cartRepository) UpdateQuantityLess(cart_id uint, quantity uint) (models.CartDetails, error) {

	var cartItem domain.CartItems

	result := ad.DB.Model(&cartItem).Where("id = ?", cart_id).Update("quantity", gorm.Expr("quantity - ?", quantity))

	if result.Error != nil {
		return models.CartDetails{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.CartDetails{}, errors.New("cart item not found")
	}

	return models.CartDetails{
		ID:        cartItem.ID,
		UserID:    cartItem.UserID,
		ProductID: cartItem.ProductID,
		Quantity:  cartItem.Quantity,
	}, nil
}

func (ad *cartRepository) UpdateQuantity(cart_id uint, quantity uint) (models.CartDetails, error) {

	var cartItem domain.CartItems

	result := ad.DB.Model(&cartItem).Where("id = ?", cart_id).Update("quantity", quantity)

	if result.Error != nil {
		return models.CartDetails{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.CartDetails{}, errors.New("cart item not found")
	}

	return models.CartDetails{
		ID:        cartItem.ID,
		UserID:    cartItem.UserID,
		ProductID: cartItem.ProductID,
		Quantity:  cartItem.Quantity,
	}, nil
}

func (ad *cartRepository) AddToCart(user_id, product_id uint, quantity uint) (models.CartDetails, error) {

	var cartDetails models.CartDetails

	err := ad.DB.Raw(`INSERT INTO cart_items (user_id,product_id,quantity) VALUES ($1,$2,$3) RETURNING id, user_id, product_id, quantity`,
		user_id, product_id, quantity).Scan(&cartDetails).Error
	if err != nil {
		return models.CartDetails{}, err
	}

	return cartDetails, nil
}
