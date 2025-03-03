package repository

import (
	"ahava/pkg/domain"
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type AdminRepository interface {
	Login(adminDetails models.AdminLogin) (domain.Admin, error)
	ListAllUsers(limit, offset int) (models.ListUsers, error)
	UpdateBlockUser(user_id uint, is_blocked bool) error

	// NewPaymentMethod(string) error
	// ListPaymentMethods() ([]domain.PaymentMethod, error)
	// CheckIfPaymentMethodAlreadyExists(payment string) (bool, error)
	// DeletePaymentMethod(id uint) error
}

type adminRepository struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) AdminRepository {
	return &adminRepository{
		DB: DB,
	}
}

func (r *adminRepository) Login(adminDetails models.AdminLogin) (domain.Admin, error) {

	var adminCompareDetails domain.Admin
	if err := r.DB.Raw("SELECT * FROM admins WHERE email = ?", adminDetails.Email).Scan(&adminCompareDetails).Error; err != nil {
		return domain.Admin{}, err
	}

	return adminCompareDetails, nil
}

// func (r *adminRepository) GetUser(user_id uint) (domain.User, error) {
// 	// Define the user
// 	var user domain.User
// 	// Query to get the user details
// 	err := r.DB.First(&user, user_id).Error
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return domain.User{}, models.ErrEntityNotFound
// 		}
// 		return domain.User{}, err
// 	}
// 	// Return the user details
// 	return user, nil
// }

func (r *adminRepository) UpdateBlockUser(user_id uint, is_blocked bool) error {
	// Update the user status
	result := r.DB.Model(&domain.User{}).
		Where("id = ?", user_id).
		Update("is_blocked", is_blocked)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrEntityNotFound
	}
	return nil
}

func (r *adminRepository) ListAllUsers(limit, offset int) (models.ListUsers, error) {
	// Define list of users and user details
	var detailts []models.UserDetailsAtAdmin
	var total int64
	// Define the query
	query := r.DB.Model(&domain.User{})
	// Query to get the total number of users
	if err := query.Count(&total).Error; err != nil {
		return models.ListUsers{}, err
	}
	// Query to get the user details
	if err := query.Offset(offset).Limit(limit).Find(&detailts).Error; err != nil {
		return models.ListUsers{}, err
	}
	// Return the list of users
	return models.ListUsers{
		Users:  detailts,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// func (r *adminRepository) NewPaymentMethod(pay string) error {

// 	if err := r.DB.Exec("insert into payment_methods(payment_name)values($1)", pay).Error; err != nil {
// 		return err
// 	}

// 	return nil

// }

// func (r *adminRepository) ListPaymentMethods() ([]domain.PaymentMethod, error) {
// 	var model []domain.PaymentMethod
// 	err := r.DB.Raw("SELECT * FROM payment_methods where is_deleted = false").Scan(&model).Error
// 	if err != nil {
// 		return []domain.PaymentMethod{}, err
// 	}

// 	return model, nil
// }

// func (r *adminRepository) CheckIfPaymentMethodAlreadyExists(payment string) (bool, error) {
// 	var count int64
// 	err := r.DB.Raw("SELECT COUNT(*) FROM payment_methods WHERE payment_name = $1 and is_deleted = false", payment).Scan(&count).Error
// 	if err != nil {
// 		return false, err
// 	}

// 	return count > 0, nil
// }

// func (r *adminRepository) DeletePaymentMethod(id uint) error {
// 	err := r.DB.Exec("UPDATE payment_methods SET is_deleted = true WHERE id = $1 ", id).Error
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
