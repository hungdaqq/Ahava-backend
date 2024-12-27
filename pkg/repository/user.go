package repository

import (
	"errors"

	"ahava/pkg/domain"
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	UserSignUp(user models.UserDetails, referal string) (models.UserDetailsResponse, error)
	CheckUserAvailability(email, phone string) bool
	FindUserByEmail(user models.UserLogin) (models.UserSignInResponse, error)
	UserBlockStatus(email string) (bool, error)
	AddAddress(user_id uint, address models.Address) (models.Address, error)
	GetAddresses(user_id uint) ([]models.Address, error)
	UpdateAddress(address_id uint, address models.Address) (models.Address, error)
	DeleteAddress(address_id uint) error

	GetUserDetails(user_id uint) (models.UserDetailsResponse, error)
	ChangePassword(user_id uint, password string) error
	GetPassword(user_id uint) (string, error)
	// FindIdFromPhone(phone string) (int, error)
	EditProfile(user_id uint, name, email, phone string) (models.UserDetailsResponse, error)

	// CheckIfFirstAddress(user_id uint) bool

	// CreditReferencePointsToWallet(user_id uint) error
	// FindUserFromReference(ref string) (int, error)

	// GetReferralCodeFromID(id uint) (string, error)
}

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) UserRepository {
	return &userDatabase{DB}
}

func (c *userDatabase) CheckUserAvailability(email, phone string) bool {

	var count int

	if err := c.DB.Raw(`SELECT COUNT(*) FROM users WHERE email = $1 OR phone = $2`,
		email, phone).Scan(&count).Error; err != nil {
		return false
	}

	return count > 0
}

func (c *userDatabase) UserSignUp(user models.UserDetails, referral string) (models.UserDetailsResponse, error) {

	var userDetails models.UserDetailsResponse
	err := c.DB.Raw("INSERT INTO users (name, email, password, phone, referral_code, birth_date) VALUES (?, ?, ?, ?, ?, ?) RETURNING *",
		user.Name, user.Email, user.Password, user.Phone, referral, user.BirthDate).Scan(&userDetails).Error

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

	err := c.DB.Raw(`SELECT * FROM users where email = ? and blocked = false`,
		user.Email).Scan(&user_details).Error
	if err != nil {
		return models.UserSignInResponse{}, errors.New("error checking user details")
	}

	return user_details, nil

}

func (i *userDatabase) AddAddress(user_id uint, address models.Address) (models.Address, error) {

	var addAddress models.Address
	if address.Type == "" {
		address.Type = "HOME"
	}
	err := i.DB.Raw(`INSERT INTO addresses (user_id, name, street, ward, district, city, phone, type, "default") 
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *`,
		user_id, address.Name, address.Street, address.Ward, address.District, address.City, address.Phone, address.Type, address.Default).Scan(&addAddress).Error
	if err != nil {
		return models.Address{}, errors.New("could not add address")
	}

	return addAddress, nil
}

func (i *userDatabase) UpdateAddress(address_id uint, address models.Address) (models.Address, error) {

	var addAddress models.Address

	err := i.DB.Raw(`UPDATE addresses SET name=$1,street=$2,ward=$3,district=$4,city=$5,phone=$6,type=$7,"default"=$8 
					WHERE id=$9 RETURNING *`,
		address.Name, address.Street, address.Ward, address.District, address.City, address.Phone, address.Type, address.Default, address_id).Scan(&addAddress).Error
	if err != nil {
		return models.Address{}, errors.New("could not add address")
	}

	return addAddress, nil
}

func (i *userDatabase) DeleteAddress(address_id uint) error {

	var addAddress models.Address

	err := i.DB.Exec(`DELETE FROM addresses WHERE id=?`,
		address_id).Scan(&addAddress).Error
	if err != nil {
		return errors.New("could not delete address")
	}

	return nil
}

// func (c *userDatabase) CheckIfFirstAddress(user_id uint) bool {

// 	var count int
// 	// query := fmt.Sprintf("select count(*) from addresses where user_id='%s'", id)
// 	if err := c.DB.Raw("SELECT COUNT(*) FROM addresses WHERE user_id=$1", user_id).Scan(&count).Error; err != nil {
// 		return false
// 	}
// 	// if count is greater than 0 that means the user already exist
// 	return count > 0

// }

func (ad *userDatabase) GetAddresses(user_id uint) ([]models.Address, error) {

	var addresses []models.Address

	if err := ad.DB.Raw("SELECT * FROM addresses WHERE user_id=?",
		user_id).Scan(&addresses).Error; err != nil {
		return []models.Address{}, errors.New("error in getting addresses")
	}

	return addresses, nil

}

func (ad *userDatabase) GetUserDetails(user_id uint) (models.UserDetailsResponse, error) {

	var details models.UserDetailsResponse
	if err := ad.DB.Raw("SELECT * FROM users where id=?",
		user_id).Scan(&details).Error; err != nil {
		return models.UserDetailsResponse{}, errors.New("could not get user details")
	}

	return details, nil

}

func (i *userDatabase) ChangePassword(id uint, password string) error {

	err := i.DB.Exec("UPDATE users SET password=$1 WHERE id=$2", password, id).Error
	if err != nil {
		return err
	}

	return nil

}

func (i *userDatabase) GetPassword(id uint) (string, error) {

	var userPassword string
	err := i.DB.Raw("select password from users where id = ?", id).Scan(&userPassword).Error
	if err != nil {
		return "", err
	}
	return userPassword, nil

}

// func (ad *userDatabase) FindIdFromPhone(phone string) (uint, error) {

// 	var id uint

// 	if err := ad.DB.Raw("select id from users where phone=?", phone).Scan(&id).Error; err != nil {
// 		return id, err
// 	}

// 	return id, nil

// }

func (i *userDatabase) EditProfile(userID uint, name, email, phone string) (models.UserDetailsResponse, error) {

	var user domain.Users

	result := i.DB.Model(&user).Updates(domain.Users{
		Name:  name,
		Email: email,
		Phone: phone,
	})

	if result.Error != nil {
		return models.UserDetailsResponse{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.UserDetailsResponse{}, errors.New("user not found")
	}

	return models.UserDetailsResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		BirthDate: user.BirthDate,
		CreateAt:  user.CreateAt,
		UpdateAt:  user.UpdateAt,
	}, nil
}

func (ad *userDatabase) GetCartID(id uint) (uint, error) {

	var cart_id uint

	if err := ad.DB.Raw("select id from carts where user_id=?", id).Scan(&cart_id).Error; err != nil {
		return 0, err
	}

	return cart_id, nil

}

// func (ad *userDatabase) FindUserFromReference(ref string) (uint, error) {
// 	var user int

// 	if err := ad.DB.Raw("SELECT id FROM users WHERE referral_code = ?", ref).Find(&user).Error; err != nil {
// 		return 0, err
// 	}

// 	return user, nil
// }

func (i *userDatabase) CreditReferencePointsToWallet(user_id uint) error {
	err := i.DB.Exec("Update wallets set amount=amount+20 where user_id=$1", user_id).Error
	if err != nil {
		return err
	}

	return nil
}

func (i *userDatabase) GetReferralCodeFromID(id uint) (string, error) {
	var referral string
	err := i.DB.Raw("SELECT referral_code FROM users WHERE id=?", id).Scan(&referral).Error
	if err != nil {
		return "", err
	}

	return referral, nil
}

func (i *userDatabase) FindProductImage(id uint) (string, error) {
	var image string
	err := i.DB.Raw("SELECT image FROM products WHERE id = ?", id).Scan(&image).Error
	if err != nil {
		return "", err
	}

	return image, nil
}

func (i *userDatabase) FindStock(id uint) (int, error) {
	var stock int
	err := i.DB.Raw("SELECT stock FROM products WHERE id = ?", id).Scan(&stock).Error
	if err != nil {
		return 0, err
	}

	return stock, nil
}
