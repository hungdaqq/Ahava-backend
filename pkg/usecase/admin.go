package usecase

import (
	"errors"

	domain "ahava/pkg/domain"
	helper "ahava/pkg/helper"
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type AdminUseCase interface {
	LoginHandler(adminDetails models.AdminLogin) (domain.TokenAdmin, error)
	BlockUser(user_id uint) error
	UnBlockUser(user_id uint) error
	GetUsers(page int) ([]models.UserDetailsAtAdmin, error)
	NewPaymentMethod(string) error
	ListPaymentMethods() ([]domain.PaymentMethod, error)
	DeletePaymentMethod(id uint) error
}

type adminUseCase struct {
	adminRepository repository.AdminRepository
	helper          helper.Helper
}

func NewAdminUseCase(repo repository.AdminRepository, h helper.Helper) AdminUseCase {
	return &adminUseCase{
		adminRepository: repo,
		helper:          h,
	}
}

func (ad *adminUseCase) LoginHandler(adminDetails models.AdminLogin) (domain.TokenAdmin, error) {

	// getting details of the admin based on the email provided
	adminCompareDetails, err := ad.adminRepository.LoginHandler(adminDetails)
	if err != nil {
		return domain.TokenAdmin{}, err
	}

	// compare password from database and that provided from admins
	err = bcrypt.CompareHashAndPassword([]byte(adminCompareDetails.Password), []byte(adminDetails.Password))
	if err != nil {
		return domain.TokenAdmin{}, err
	}

	var adminDetailsResponse models.AdminDetailsResponse

	//  copy all details except password and sent it back to the front end
	err = copier.Copy(&adminDetailsResponse, &adminCompareDetails)
	if err != nil {
		return domain.TokenAdmin{}, err
	}

	access, refresh, err := ad.helper.GenerateTokenAdmin(adminDetailsResponse)

	if err != nil {
		return domain.TokenAdmin{}, err
	}

	return domain.TokenAdmin{
		Admin:        adminDetailsResponse,
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil

}

func (ad *adminUseCase) BlockUser(user_id uint) error {

	user, err := ad.adminRepository.GetUserByID(user_id)
	if err != nil {
		return err
	}

	if user.IsBlocked {
		return errors.New("already blocked")
	} else {
		user.IsBlocked = true
	}

	err = ad.adminRepository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}

	return nil

}

// unblock user
func (ad *adminUseCase) UnBlockUser(user_id uint) error {

	user, err := ad.adminRepository.GetUserByID(user_id)
	if err != nil {
		return err
	}

	if user.IsBlocked {
		user.IsBlocked = false
	} else {
		return errors.New("already unblocked")
	}

	err = ad.adminRepository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}

	return nil

}

func (ad *adminUseCase) GetUsers(page int) ([]models.UserDetailsAtAdmin, error) {

	userDetails, err := ad.adminRepository.GetUsers(page)
	if err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}

	return userDetails, nil

}

func (ad *adminUseCase) NewPaymentMethod(id string) error {

	exists, err := ad.adminRepository.CheckIfPaymentMethodAlreadyExists(id)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("payment method already exists")
	}

	err = ad.adminRepository.NewPaymentMethod(id)
	if err != nil {
		return err
	}

	return nil

}

func (ad *adminUseCase) ListPaymentMethods() ([]domain.PaymentMethod, error) {

	categories, err := ad.adminRepository.ListPaymentMethods()
	if err != nil {
		return []domain.PaymentMethod{}, err
	}
	return categories, nil

}

func (ad *adminUseCase) DeletePaymentMethod(id uint) error {

	err := ad.adminRepository.DeletePaymentMethod(id)
	if err != nil {
		return err
	}
	return nil

}
