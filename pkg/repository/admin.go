package repository

import (
	errors "ahava/pkg/utils/errors"
	"fmt"

	"ahava/pkg/domain"
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type AdminRepository interface {
	LoginHandler(adminDetails models.AdminLogin) (domain.Admin, error)
	GetUserByID(user_id uint) (domain.User, error)
	UpdateBlockUserByID(user domain.User) error
	GetUsers(page int) ([]models.UserDetailsAtAdmin, error)
	NewPaymentMethod(string) error
	ListPaymentMethods() ([]domain.PaymentMethod, error)
	CheckIfPaymentMethodAlreadyExists(payment string) (bool, error)
	DeletePaymentMethod(id uint) error
}

type adminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) AdminRepository {
	return &adminRepository{
		DB: DB,
	}
}

func (r *adminRepository) LoginHandler(adminDetails models.AdminLogin) (domain.Admin, error) {

	var adminCompareDetails domain.Admin
	if err := r.DB.Raw("SELECT * FROM admins WHERE email = ?", adminDetails.Email).Scan(&adminCompareDetails).Error; err != nil {
		return domain.Admin{}, err
	}

	return adminCompareDetails, nil
}

func (r *adminRepository) GetUserByID(user_id uint) (domain.User, error) {
	var count int
	if err := r.DB.Raw("select count(*) from users where id = ?", user_id).Scan(&count).Error; err != nil {
		return domain.User{}, err
	}
	if count < 1 {
		return domain.User{}, errors.ErrEntityNotFound
	}

	query := fmt.Sprintf("select * from users where id = '%d'", user_id)
	var userDetails domain.User

	if err := r.DB.Raw(query).Scan(&userDetails).Error; err != nil {
		return domain.User{}, err
	}

	return userDetails, nil
}

// function which will both block and unblock a user
func (r *adminRepository) UpdateBlockUserByID(user domain.User) error {

	err := r.DB.Exec("update users set is_blocked = ? where id = ?", user.IsBlocked, user.ID).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *adminRepository) GetUsers(page int) ([]models.UserDetailsAtAdmin, error) {
	// pagination purpose -
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * 5
	var userDetails []models.UserDetailsAtAdmin

	if err := r.DB.Raw("select id,name,email,phone,blocked from users limit ? offset ?", 20, offset).Scan(&userDetails).Error; err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}

	return userDetails, nil

}

func (r *adminRepository) NewPaymentMethod(pay string) error {

	if err := r.DB.Exec("insert into payment_methods(payment_name)values($1)", pay).Error; err != nil {
		return err
	}

	return nil

}

func (r *adminRepository) ListPaymentMethods() ([]domain.PaymentMethod, error) {
	var model []domain.PaymentMethod
	err := r.DB.Raw("SELECT * FROM payment_methods where is_deleted = false").Scan(&model).Error
	if err != nil {
		return []domain.PaymentMethod{}, err
	}

	return model, nil
}

func (r *adminRepository) CheckIfPaymentMethodAlreadyExists(payment string) (bool, error) {
	var count int64
	err := r.DB.Raw("SELECT COUNT(*) FROM payment_methods WHERE payment_name = $1 and is_deleted = false", payment).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *adminRepository) DeletePaymentMethod(id uint) error {
	err := r.DB.Exec("UPDATE payment_methods SET is_deleted = true WHERE id = $1 ", id).Error
	if err != nil {
		return err
	}

	return nil
}
