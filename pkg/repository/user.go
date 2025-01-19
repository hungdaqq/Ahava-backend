package repository

import (
	"ahava/pkg/domain"
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Register(user models.UserDetails, referal string) (models.UserDetailsResponse, error)
	CheckUserAvailability(email, username string) bool
	FindUser(user models.UserLogin) (domain.User, error)
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

func (r *userDatabase) Register(user models.UserDetails, referral string) (models.UserDetailsResponse, error) {

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

func (r *userDatabase) FindUser(login models.UserLogin) (domain.User, error) {

	var user domain.User

	err := r.DB.First(&user, "email = ? OR username = ?", login.Email, login.Username).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.User{}, models.ErrEntityNotFound
		}
		return domain.User{}, err
	}
	if user.IsBlocked {
		return domain.User{}, models.ErrForbidden
	}

	return user, nil
}

func (r *userDatabase) AddAddress(user_id uint, a models.Address) (models.Address, error) {

	address := domain.Address{
		UserID:       user_id,
		Name:         a.Name,
		Street:       a.Street,
		Ward:         a.Ward,
		WardCode:     a.WardCode,
		District:     a.District,
		DistrictCode: a.DistrictCode,
		Province:     a.Province,
		ProvinceCode: a.ProvinceCode,
		Phone:        a.Phone,
		Type:         a.Type,
		Default:      a.Default,
	}

	err := r.DB.Create(&address).Error
	if err != nil {
		return models.Address{}, err
	}

	return models.Address{
		ID:           address.ID,
		Name:         address.Name,
		Street:       address.Street,
		Ward:         address.Ward,
		WardCode:     address.WardCode,
		District:     address.District,
		DistrictCode: address.DistrictCode,
		Province:     address.Province,
		ProvinceCode: address.ProvinceCode,
		Phone:        address.Phone,
		Type:         address.Type,
		Default:      address.Default,
	}, nil
}

func (r *userDatabase) UpdateAddress(user_id, address_id uint, a models.Address) (models.Address, error) {

	var address models.Address

	result := r.DB.
		Model(&domain.Address{}).
		Where("id = ? AND user_id =?", address_id, user_id).
		Updates(domain.Address{
			Name:         a.Name,
			Street:       a.Street,
			Ward:         a.Ward,
			WardCode:     a.WardCode,
			District:     a.District,
			DistrictCode: a.DistrictCode,
			Province:     a.Province,
			ProvinceCode: a.ProvinceCode,
			Phone:        a.Phone,
			Type:         a.Type,
			Default:      a.Default,
		}).
		Scan(&address)

	if result.Error != nil {
		return models.Address{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.Address{}, models.ErrEntityNotFound
	}

	return address, nil
}

func (r *userDatabase) DeleteAddress(user_id, address_id uint) error {

	result := r.DB.Where("user_id = ? AND id = ?", user_id, address_id).Delete(&domain.Address{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrEntityNotFound
	}

	return nil
}

func (ad *userDatabase) GetAddresses(user_id uint) ([]models.Address, error) {

	var adresses []models.Address

	err := ad.DB.Model(&domain.Address{}).
		Where("user_id = ?", user_id).
		Find(&adresses).Error
	if err != nil {
		return []models.Address{}, err
	}

	return adresses, nil
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
		return models.ErrEntityNotFound
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

	result := r.DB.Model(&domain.User{}).
		Where("id = ?", userID).
		Updates(domain.User{
			Name:      profile.Name,
			Phone:     profile.Phone,
			BirthDate: profile.BirthDate,
			Gender:    profile.Gender,
		}).
		Scan(&user)

	if result.Error != nil {
		return models.UserDetailsResponse{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.UserDetailsResponse{}, models.ErrEntityNotFound
	}

	return user, nil
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
