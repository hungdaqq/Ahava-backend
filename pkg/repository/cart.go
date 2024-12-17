package repository

import (
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	GetCart(user_id int) ([]models.CartItem, error)
	AddToCart(user_id, product_id int, quantity uint) (models.CartDetails, error)
	CheckIfItemIsAlreadyAdded(user_id, product_id int) (int, error)
	UpdateQuantityAdd(cart_id int, quantity uint) (models.CartDetails, error)
	UpdateQuantityLess(cart_id int, quantity uint) (models.CartDetails, error)
	UpdateQuantity(cart_id int, quantity uint) (models.CartDetails, error)

	RemoveFromCart(cart_id int) error

	// GetAddresses(id int) ([]models.Address, error)
	// GetPaymentOptions() ([]models.PaymentMethod, error)
	// GetCartId(user_id int) (int, error)
	// CreateNewCart(user_id int) (int, error)

	// AddLineItems(cart_id, product_id int) error
}

type cartRepository struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) *cartRepository {
	return &cartRepository{
		DB: db,
	}
}

// func (ad *cartRepository) GetAddresses(id int) ([]models.Address, error) {

// 	var addresses []models.Address

// 	if err := ad.DB.Raw("SELECT * FROM addresses WHERE user_id=$1", id).Scan(&addresses).Error; err != nil {
// 		return []models.Address{}, err
// 	}

// 	return addresses, nil

// }

func (ad *cartRepository) GetCart(user_id int) ([]models.CartItem, error) {

	var cart []models.CartItem

	if err := ad.DB.Raw(`SELECT ci.id AS cart_id, p.id as product_id, p.product_name, p.image, p.price, ci.quantity, (ci.quantity * p.price) AS item_price 
						FROM cart_items ci JOIN products p ON ci.product_id = p.id WHERE ci.user_id = $1;`, user_id).Scan(&cart).Error; err != nil {
		return []models.CartItem{}, err
	}

	return cart, nil
}

func (ad *cartRepository) GetPaymentOptions() ([]models.PaymentMethod, error) {

	var payment []models.PaymentMethod

	if err := ad.DB.Raw("SELECT * FROM payment_methods WHERE is_deleted = false").Scan(&payment).Error; err != nil {
		return []models.PaymentMethod{}, err
	}

	return payment, nil

}

// func (ad *cartRepository) GetCartId(user_id int) (int, error) {

// 	var id int

// 	if err := ad.DB.Raw("SELECT id FROM carts WHERE user_id=?", user_id).Scan(&id).Error; err != nil {
// 		return 0, err
// 	}

// 	return id, nil

// }

// func (i *cartRepository) CreateNewCart(user_id int) (int, error) {
// 	var id int
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

// func (i *cartRepository) AddLineItems(cart_id, product_id int) error {

// 	err := i.DB.Exec(`
// 		INSERT INTO line_items (cart_id,product_id)
// 		VALUES ($1,$2)`, cart_id, product_id).Error
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (ad *cartRepository) CheckIfItemIsAlreadyAdded(user_id, product_id int) (int, error) {

	var count int

	if err := ad.DB.Raw("SELECT id FROM cart_items WHERE user_id = $1 AND product_id=$2",
		user_id, product_id).Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (ad *cartRepository) RemoveFromCart(cart_id int) error {

	err := ad.DB.Exec(`DELETE FROM cart_items WHERE id=$1`,
		cart_id).Error
	if err != nil {
		return err
	}

	return nil

}

func (ad *cartRepository) UpdateQuantityAdd(cart_id int, quantity uint) (models.CartDetails, error) {

	var cartDetails models.CartDetails
	err := ad.DB.Raw(`UPDATE cart_items SET quantity = quantity+$1 WHERE id=$2 RETURNING id, user_id, product_id, quantity`,
		quantity, cart_id).Scan(&cartDetails).Error
	if err != nil {
		return models.CartDetails{}, err
	}

	return cartDetails, nil
}

func (ad *cartRepository) UpdateQuantityLess(cart_id int, quantity uint) (models.CartDetails, error) {

	var cartDetails models.CartDetails

	err := ad.DB.Raw(`UPDATE cart_items SET quantity = quantity-$1 WHERE id=$2 RETURNING id, user_id, product_id, quantity`,
		quantity, cart_id).Scan(&cartDetails).Error
	if err != nil {
		return models.CartDetails{}, err
	}

	return cartDetails, nil
}

func (ad *cartRepository) UpdateQuantity(cart_id int, quantity uint) (models.CartDetails, error) {

	var cartDetails models.CartDetails

	err := ad.DB.Raw(`UPDATE cart_items SET quantity=$1 WHERE id=$2 RETURNING id, user_id, product_id, quantity`,
		quantity, cart_id).Scan(&cartDetails).Error
	if err != nil {
		return models.CartDetails{}, err
	}

	return cartDetails, nil
}

func (ad *cartRepository) AddToCart(user_id, product_id int, quantity uint) (models.CartDetails, error) {

	var cartDetails models.CartDetails

	err := ad.DB.Raw(`INSERT INTO cart_items (user_id,product_id,quantity) VALUES ($1,$2,$3) RETURNING id, user_id, product_id, quantity`,
		user_id, product_id, quantity).Scan(&cartDetails).Error
	if err != nil {
		return models.CartDetails{}, err
	}

	return cartDetails, nil
}
