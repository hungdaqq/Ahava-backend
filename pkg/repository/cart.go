package repository

import (
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	GetCart(user_id int) ([]models.CartItem, error)
	UpdateQuantityAdd(user_id, product_id, quantity int) error
	UpdateQuantityLess(id, inv_id int) error
	AddToCart(user_id, product_id, quantity int) error

	GetAddresses(id int) ([]models.Address, error)
	GetPaymentOptions() ([]models.PaymentMethod, error)
	GetCartId(user_id int) (int, error)
	CreateNewCart(user_id int) (int, error)

	// AddLineItems(cart_id, product_id int) error
	CheckIfItemIsAlreadyAdded(user_id, product_id int) (bool, error)

	RemoveFromCart(cart, product int) error
}

type cartRepository struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) *cartRepository {
	return &cartRepository{
		DB: db,
	}
}

func (ad *cartRepository) GetAddresses(id int) ([]models.Address, error) {

	var addresses []models.Address

	if err := ad.DB.Raw("SELECT * FROM addresses WHERE user_id=$1", id).Scan(&addresses).Error; err != nil {
		return []models.Address{}, err
	}

	return addresses, nil

}

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

func (ad *cartRepository) GetCartId(user_id int) (int, error) {

	var id int

	if err := ad.DB.Raw("SELECT id FROM carts WHERE user_id=?", user_id).Scan(&id).Error; err != nil {
		return 0, err
	}

	return id, nil

}

func (i *cartRepository) CreateNewCart(user_id int) (int, error) {
	var id int
	err := i.DB.Exec(`
		INSERT INTO carts (user_id)
		VALUES ($1)`, user_id).Error
	if err != nil {
		return 0, err
	}

	if err := i.DB.Raw("select id from carts where user_id=?", user_id).Scan(&id).Error; err != nil {
		return 0, err
	}

	return id, nil
}

func (i *cartRepository) AddLineItems(cart_id, product_id int) error {

	err := i.DB.Exec(`
		INSERT INTO line_items (cart_id,product_id)
		VALUES ($1,$2)`, cart_id, product_id).Error
	if err != nil {
		return err
	}

	return nil
}

func (ad *cartRepository) CheckIfItemIsAlreadyAdded(user_id, product_id int) (bool, error) {

	var count int

	if err := ad.DB.Raw("SELECT COUNT(*) FROM cart_items WHERE user_id = $1 AND product_id = $2",
		user_id, product_id).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (ad *cartRepository) RemoveFromCart(cart, product int) error {

	if err := ad.DB.Exec(`DELETE FROM line_items WHERE cart_id = $1 AND product_id = $2`, cart, product).Error; err != nil {
		return err
	}

	return nil

}

func (ad *cartRepository) UpdateQuantityAdd(user_id, product_id, quantity int) error {

	err := ad.DB.Exec(`UPDATE cart_items SET quantity = quantity+$1 WHERE user_id=$2 AND product_id=$3`,
		quantity, user_id, product_id).Error
	if err != nil {
		return err
	}

	return nil
}

func (ad *cartRepository) UpdateQuantityLess(id, inv_id int) error {

	if err := ad.DB.Exec(`UPDATE line_items SET quantity = quantity - 1 WHERE cart_id = $1 AND product_id=$2;`,
		id, inv_id).Error; err != nil {
		return err
	}

	return nil
}

func (ad *cartRepository) AddToCart(user_id, product_id, quantity int) error {

	err := ad.DB.Exec(`INSERT INTO cart_items (user_id,product_id,quantity) VALUES ($1,$2,$3)`,
		user_id, product_id, quantity).Error
	if err != nil {
		return err
	}

	return nil
}
