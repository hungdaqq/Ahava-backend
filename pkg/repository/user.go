package repository

import (
	"errors"
	"fmt"

	"ahava/pkg/domain"
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	UserSignUp(user models.UserDetails, referal string) (models.UserDetailsResponse, error)
	CheckUserAvailability(email string) bool
	FindUserByEmail(user models.UserLogin) (models.UserSignInResponse, error)
	UserBlockStatus(email string) (bool, error)
	AddAddress(id int, address models.AddAddress, result bool) error
	GetAddresses(id int) ([]domain.Address, error)
	GetUserDetails(id int) (models.UserDetailsResponse, error)
	ChangePassword(id int, password string) error
	GetPassword(id int) (string, error)
	FindIdFromPhone(phone string) (int, error)
	EditProfile(user_id int, name, email, phone string) (models.UserDetailsResponse, error)


	CheckIfFirstAddress(id int) bool

	GetCartID(id int) (int, error)
	GetProductsInCart(cart_id int) ([]int, error)
	FindProductNames(product_id int) (string, error)
	FindCartQuantity(cart_id, product_id int) (int, error)
	FindPrice(product_id int) (float64, error)
	FindCategory(product_id int) (int, error)
	FindofferPercentage(category_id int) (int, error)

	CreditReferencePointsToWallet(user_id int) error
	FindUserFromReference(ref string) (int, error)

	GetReferralCodeFromID(id int) (string, error)

	FindProductImage(id int) (string, error)
	FindStock(id int) (int, error)
}

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) UserRepository {
	return &userDatabase{DB}
}

func (c *userDatabase) CheckUserAvailability(email string) bool {

	var count int
	query := fmt.Sprintf("select count(*) from users where email='%s'", email)
	if err := c.DB.Raw(query).Scan(&count).Error; err != nil {
		return false
	}
	// if count is greater than 0 that means the user already exist
	return count > 0

}

func (c *userDatabase) UserSignUp(user models.UserDetails, referral string) (models.UserDetailsResponse, error) {

	var userDetails models.UserDetailsResponse
	err := c.DB.Raw("INSERT INTO users (name, email, password, phone,referral_code, birth_date, address) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id, name, email, phone",
		user.Name, user.Email, user.Password, user.Phone, referral, user.BirthDate, user.Address).Scan(&userDetails).Error

	if err != nil {
		return models.UserDetailsResponse{}, err
	}

	return userDetails, nil
}

func (cr *userDatabase) UserBlockStatus(email string) (bool, error) {
	var isBlocked bool
	err := cr.DB.Raw("select blocked from users where email = ?", email).Scan(&isBlocked).Error
	if err != nil {
		return false, err
	}
	return isBlocked, nil
}

func (c *userDatabase) FindUserByEmail(user models.UserLogin) (models.UserSignInResponse, error) {

	var user_details models.UserSignInResponse

	err := c.DB.Raw(`
		SELECT *
		FROM users where email = ? and blocked = false
		`, user.Email).Scan(&user_details).Error

	if err != nil {
		return models.UserSignInResponse{}, errors.New("error checking user details")
	}

	return user_details, nil

}

func (i *userDatabase) AddAddress(id int, address models.AddAddress, result bool) error {
	err := i.DB.Exec(`
		INSERT INTO addresses (user_id, name, house_name, street, city, state, phone, pin,"default")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9 )`,
		id, address.Name, address.HouseName, address.Street, address.City, address.State, address.Phone, address.Pin, result).Error
	if err != nil {
		return errors.New("could not add address")
	}

	return nil
}

func (c *userDatabase) CheckIfFirstAddress(id int) bool {

	var count int
	// query := fmt.Sprintf("select count(*) from addresses where user_id='%s'", id)
	if err := c.DB.Raw("select count(*) from addresses where user_id=$1", id).Scan(&count).Error; err != nil {
		return false
	}
	// if count is greater than 0 that means the user already exist
	return count > 0

}

func (ad *userDatabase) GetAddresses(user_id int) ([]domain.Address, error) {

	var addresses []domain.Address

	if err := ad.DB.Raw("SELECT * FROM addresses WHERE user_id=?",
		user_id).Scan(&addresses).Error; err != nil {
		return []domain.Address{}, errors.New("error in getting addresses")
	}

	return addresses, nil

}

func (ad *userDatabase) GetUserDetails(user_id int) (models.UserDetailsResponse, error) {

	var details models.UserDetailsResponse
	if err := ad.DB.Raw("SELECT id,name,email,phone,birth_date,create_at FROM users where id=?",
		user_id).Scan(&details).Error; err != nil {
		return models.UserDetailsResponse{}, errors.New("could not get user details")
	}

	return details, nil

}

func (i *userDatabase) ChangePassword(id int, password string) error {

	err := i.DB.Exec("UPDATE users SET password=$1 WHERE id=$2", password, id).Error
	if err != nil {
		return err
	}

	return nil

}

func (i *userDatabase) GetPassword(id int) (string, error) {

	var userPassword string
	err := i.DB.Raw("select password from users where id = ?", id).Scan(&userPassword).Error
	if err != nil {
		return "", err
	}
	return userPassword, nil

}

func (ad *userDatabase) FindIdFromPhone(phone string) (int, error) {

	var id int

	if err := ad.DB.Raw("select id from users where phone=?", phone).Scan(&id).Error; err != nil {
		return id, err
	}

	return id, nil

}

func (i *userDatabase) EditProfile(user_id int, name, email, phone string) (models.UserDetailsResponse, error) {

	var user models.UserDetailsResponse

	err := i.DB.Raw(`UPDATE users SET name=$1, phone=$2 WHERE id=$3 RETURNING *`,
		name, phone, user_id).Scan(&user).Error
	if err != nil {
		return models.UserDetailsResponse{}, err
	}

	return user, nil
}

func (ad *userDatabase) GetCartID(id int) (int, error) {

	var cart_id int

	if err := ad.DB.Raw("select id from carts where user_id=?", id).Scan(&cart_id).Error; err != nil {
		return 0, err
	}

	return cart_id, nil

}

func (ad *userDatabase) GetProductsInCart(cart_id int) ([]int, error) {

	var cart_products []int

	if err := ad.DB.Raw("select product_id from line_items where cart_id=?", cart_id).Scan(&cart_products).Error; err != nil {
		return []int{}, err
	}

	return cart_products, nil

}

func (ad *userDatabase) FindProductNames(product_id int) (string, error) {

	var product_name string

	if err := ad.DB.Raw("select product_name from products where id=?", product_id).Scan(&product_name).Error; err != nil {
		return "", err
	}

	return product_name, nil

}

func (ad *userDatabase) FindCartQuantity(cart_id, product_id int) (int, error) {

	var quantity int

	if err := ad.DB.Raw("select quantity from line_items where cart_id=$1 and product_id=$2", cart_id, product_id).Scan(&quantity).Error; err != nil {
		return 0, err
	}

	return quantity, nil

}

func (ad *userDatabase) FindPrice(product_id int) (float64, error) {

	var price float64

	if err := ad.DB.Raw("select price from products where id=?", product_id).Scan(&price).Error; err != nil {
		return 0, err
	}

	return price, nil

}

func (ad *userDatabase) FindCategory(product_id int) (int, error) {

	var category int

	if err := ad.DB.Raw("select category_id from products where id=?", product_id).Scan(&category).Error; err != nil {
		return 0, err
	}

	return category, nil

}

func (ad *userDatabase) FindofferPercentage(category_id int) (int, error) {
	var percentage int
	err := ad.DB.Raw("select discount_rate from offers where category_id=$1 and valid=true", category_id).Scan(&percentage).Error
	if err != nil {
		return 0, err
	}

	return percentage, nil
}

func (ad *userDatabase) FindUserFromReference(ref string) (int, error) {
	var user int

	if err := ad.DB.Raw("SELECT id FROM users WHERE referral_code = ?", ref).Find(&user).Error; err != nil {
		return 0, err
	}

	return user, nil
}

func (i *userDatabase) CreditReferencePointsToWallet(user_id int) error {
	err := i.DB.Exec("Update wallets set amount=amount+20 where user_id=$1", user_id).Error
	if err != nil {
		return err
	}

	return nil
}

func (i *userDatabase) GetReferralCodeFromID(id int) (string, error) {
	var referral string
	err := i.DB.Raw("SELECT referral_code FROM users WHERE id=?", id).Scan(&referral).Error
	if err != nil {
		return "", err
	}

	return referral, nil
}

func (i *userDatabase) FindProductImage(id int) (string, error) {
	var image string
	err := i.DB.Raw("SELECT image FROM products WHERE id = ?", id).Scan(&image).Error
	if err != nil {
		return "", err
	}

	return image, nil
}

func (i *userDatabase) FindStock(id int) (int, error) {
	var stock int
	err := i.DB.Raw("SELECT stock FROM products WHERE id = ?", id).Scan(&stock).Error
	if err != nil {
		return 0, err
	}

	return stock, nil
}
