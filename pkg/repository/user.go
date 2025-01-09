package repository

import (
	errors "ahava/pkg/utils/errors"

	"ahava/pkg/domain"
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	UserSignUp(user models.UserDetails, referal string) (models.UserDetailsResponse, error)
	CheckUserAvailability(email, username string) bool
	FindUser(user models.UserLogin) (models.UserSignInResponse, error)
	UserBlockStatus(email, username string) (bool, error)
	AddAddress(user_id uint, address models.Address) (models.Address, error)
	GetAddresses(user_id uint) ([]models.Address, error)
	UpdateAddress(user_id, address_id uint, address models.Address) (models.Address, error)
	DeleteAddress(user_id, address_id uint) error

	GetUserDetails(user_id uint) (models.UserDetailsResponse, error)
	ChangePassword(user_id uint, password string) error
	GetPassword(user_id uint) (string, error)
	// FindIdFromPhone(phone string) (int, error)
	EditProfile(user_id uint, profile models.EditProfile) (models.UserDetailsResponse, error)

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

func (r *userDatabase) CheckUserAvailability(email, phone string) bool {

	var count int

	if err := r.DB.Raw(`SELECT COUNT(*) FROM users WHERE email = $1 OR username = $2`,
		email, phone).Scan(&count).Error; err != nil {
		return false
	}

	return count > 0
}

func (r *userDatabase) UserSignUp(user models.UserDetails, referral string) (models.UserDetailsResponse, error) {

	var userDetails models.UserDetailsResponse

	err := r.DB.
		Create(&domain.User{
			Name:         user.Name,
			Username:     user.Username,
			Gender:       user.Gender,
			Email:        user.Email,
			Password:     user.Password,
			Phone:        user.Phone,
			ReferralCode: referral,
			BirthDate:    user.BirthDate,
		}).
		Scan(&userDetails).Error
	if err != nil {
		return models.UserDetailsResponse{}, err
	}

	return userDetails, nil
}

func (r *userDatabase) UserBlockStatus(email, username string) (bool, error) {
	var isBlocked bool
	err := r.DB.Raw("SELECT is_blocked FROM users WHERE email=? OR username=?",
		email, username).Scan(&isBlocked).Error
	if err != nil {
		return false, err
	}
	return isBlocked, nil
}

func (r *userDatabase) FindUser(login models.UserLogin) (models.UserSignInResponse, error) {

	var response models.UserSignInResponse

	err := r.DB.Raw(`SELECT * FROM users WHERE (email = ? OR username = ?) AND is_blocked = false`,
		login.Email, login.Username).Scan(&response).Error
	if err != nil {
		return models.UserSignInResponse{}, err
	}

	return response, nil

}

func (r *userDatabase) AddAddress(user_id uint, address models.Address) (models.Address, error) {

	var addAddress models.Address

	err := r.DB.
		Model(&domain.Address{}).
		Create(&domain.Address{
			UserID:   user_id,
			Name:     address.Name,
			Street:   address.Street,
			Ward:     address.Ward,
			District: address.District,
			Province: address.Province,
			Phone:    address.Phone,
			Type:     address.Type,
			Default:  address.Default,
		}).
		Scan(&addAddress).Error
	if err != nil {
		return models.Address{}, err
	}

	return addAddress, nil
}

func (r *userDatabase) UpdateAddress(user_id, address_id uint, address models.Address) (models.Address, error) {

	var updateAddress models.Address

	result := r.DB.
		Model(&domain.Address{}).
		Where("id = ? AND user_id =?", address_id, user_id).
		Updates(domain.Address{
			Name:     address.Name,
			Street:   address.Street,
			Ward:     address.Ward,
			District: address.District,
			Province: address.Province,
			Phone:    address.Phone,
			Type:     address.Type,
			Default:  address.Default,
		}).
		Scan(&updateAddress)

	if result.Error != nil {
		return models.Address{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.Address{}, errors.ErrEntityNotFound
	}

	return updateAddress, nil
}

func (r *userDatabase) DeleteAddress(address_id, user_id uint) error {

	result := r.DB.Exec(`DELETE FROM addresses WHERE id=? AND user_id=?`,
		address_id, user_id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrEntityNotFound
	}

	return nil
}

// func (r *userDatabase) CheckIfFirstAddress(user_id uint) bool {

// 	var count int
// 	// query := fmt.Sprintf("select count(*) from addresses where user_id='%s'", id)
// 	if err := r.DB.Raw("SELECT COUNT(*) FROM addresses WHERE user_id=$1", user_id).Scan(&count).Error; err != nil {
// 		return false
// 	}
// 	// if count is greater than 0 that means the user already exist
// 	return count > 0

// }

func (ad *userDatabase) GetAddresses(user_id uint) ([]models.Address, error) {

	var addresses []models.Address

	if err := ad.DB.Raw("SELECT * FROM addresses WHERE user_id=?",
		user_id).Scan(&addresses).Error; err != nil {
		return []models.Address{}, err
	}

	return addresses, nil

}

func (ad *userDatabase) GetUserDetails(user_id uint) (models.UserDetailsResponse, error) {

	var details models.UserDetailsResponse
	if err := ad.DB.Raw("SELECT * FROM users where id=?",
		user_id).Scan(&details).Error; err != nil {
		return models.UserDetailsResponse{}, err
	}

	return details, nil

}

func (r *userDatabase) ChangePassword(id uint, password string) error {

	result := r.DB.Model(&domain.User{}).Where("id = ?", id).Update("password", password)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.ErrEntityNotFound
	}

	return nil
}

func (r *userDatabase) GetPassword(id uint) (string, error) {

	var userPassword string
	err := r.DB.Raw("select password from users where id = ?", id).Scan(&userPassword).Error
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

func (r *userDatabase) EditProfile(userID uint, profile models.EditProfile) (models.UserDetailsResponse, error) {

	var user models.UserDetailsResponse

	result := r.DB.Model(&domain.User{}).Where("id = ?", userID).Updates(domain.User{
		Name:      profile.Name,
		Phone:     profile.Phone,
		BirthDate: profile.BirthDate,
		Gender:    profile.Gender,
	}).Scan(&user)

	if result.Error != nil {
		return models.UserDetailsResponse{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.UserDetailsResponse{}, errors.ErrEntityNotFound
	}

	return user, nil
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

func (r *userDatabase) CreditReferencePointsToWallet(user_id uint) error {
	err := r.DB.Exec("Update wallets set amount=amount+20 where user_id=$1", user_id).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *userDatabase) GetReferralCodeFromID(id uint) (string, error) {
	var referral string
	err := r.DB.Raw("SELECT referral_code FROM users WHERE id=?", id).Scan(&referral).Error
	if err != nil {
		return "", err
	}

	return referral, nil
}

func (r *userDatabase) FindProductImage(id uint) (string, error) {
	var image string
	err := r.DB.Raw("SELECT default_image FROM products WHERE id = ?", id).Scan(&image).Error
	if err != nil {
		return "", err
	}

	return image, nil
}

func (r *userDatabase) FindStock(id uint) (int, error) {
	var stock int
	err := r.DB.Raw("SELECT stock FROM products WHERE id = ?", id).Scan(&stock).Error
	if err != nil {
		return 0, err
	}

	return stock, nil
}
