package service

import (
	domain "ahava/pkg/domain"
	helper "ahava/pkg/helper"
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type AdminService interface {
	Login(models.AdminLogin) (domain.TokenAdmin, error)
	BlockUser(user_id uint) error
	UnBlockUser(user_id uint) error
	GetAllUsers(limit, offset int) (models.ListUsers, error)
	// NewPaymentMethod(string) error
	// ListPaymentMethods() ([]domain.PaymentMethod, error)
	// DeletePaymentMethod(id uint) error
}

type adminService struct {
	adminRepository repository.AdminRepository
	helper          helper.Helper
}

func NewAdminService(repo repository.AdminRepository, h helper.Helper) AdminService {
	return &adminService{
		adminRepository: repo,
		helper:          h,
	}
}

func (ad *adminService) Login(adminDetails models.AdminLogin) (domain.TokenAdmin, error) {
	// Get the admin details
	adminCompareDetails, err := ad.adminRepository.Login(adminDetails)
	if err != nil {
		return domain.TokenAdmin{}, err
	}
	// Compare password from database and that provided from admins
	err = bcrypt.CompareHashAndPassword([]byte(adminCompareDetails.Password), []byte(adminDetails.Password))
	if err != nil {
		return domain.TokenAdmin{}, err
	}
	// Copy all details except password and sent it back to the front end
	var adminDetailsResponse models.AdminDetailsResponse
	err = copier.Copy(&adminDetailsResponse, &adminCompareDetails)
	if err != nil {
		return domain.TokenAdmin{}, err
	}
	// Generate token for the admin
	access, refresh, err := ad.helper.GenerateTokenAdmin(adminDetailsResponse)
	if err != nil {
		return domain.TokenAdmin{}, err
	}
	// Return the token
	return domain.TokenAdmin{
		Admin:        adminDetailsResponse,
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (ad *adminService) BlockUser(user_id uint) error {
	// Block the user
	err := ad.adminRepository.UpdateBlockUser(user_id, true)
	if err != nil {
		return err
	}
	return nil
}

func (ad *adminService) UnBlockUser(user_id uint) error {
	// Unblock the user
	err := ad.adminRepository.UpdateBlockUser(user_id, false)
	if err != nil {
		return err
	}
	return nil
}

func (ad *adminService) GetAllUsers(limit, offset int) (models.ListUsers, error) {
	// Get all users with limit and offset
	listUsers, err := ad.adminRepository.GetAllUsers(limit, offset)
	if err != nil {
		return models.ListUsers{}, err
	}
	// Return the list of users
	return listUsers, nil
}

// func (ad *adminService) NewPaymentMethod(id string) error {

// 	exists, err := ad.adminRepository.CheckIfPaymentMethodAlreadyExists(id)
// 	if err != nil {
// 		return err
// 	}

// 	if exists {
// 		return errors.New("payment method already exists")
// 	}

// 	err = ad.adminRepository.NewPaymentMethod(id)
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }

// func (ad *adminService) ListPaymentMethods() ([]domain.PaymentMethod, error) {

// 	categories, err := ad.adminRepository.ListPaymentMethods()
// 	if err != nil {
// 		return []domain.PaymentMethod{}, err
// 	}
// 	return categories, nil

// }

// func (ad *adminService) DeletePaymentMethod(id uint) error {

// 	err := ad.adminRepository.DeletePaymentMethod(id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil

// }
