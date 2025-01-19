package service

import (
	"errors"

	"ahava/pkg/config"
	helper "ahava/pkg/helper"
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
)

type UserService interface {
	Register(user models.UserDetails, ref string) (models.TokenUsers, error)
	Login(user models.UserLogin) (models.TokenUsers, error)
	AddAddress(user_id uint, address models.Address) (models.Address, error)
	GetAddresses(user_id uint) ([]models.Address, error)
	UpdateAddress(user_id, address_id uint, address models.Address) (models.Address, error)
	DeleteAddress(user_id, address_id uint) error

	GetUserDetails(user_id uint) (models.UserDetailsResponse, error)

	ChangePassword(user_id uint, old string, password string, repassword string) error
	// ForgotPasswordSend(phone string) error
	// ForgotPasswordVerifyAndChange(model models.ForgotVerify) error

	EditProfile(user_id uint, profile models.EditProfile) (models.UserDetailsResponse, error)

	// GetMyReferenceLink(id uint) (string, error)
}

type userService struct {
	userRepo repository.UserRepository
	cfg      config.Config
	// otpRepository     repository.OtpRepository
	// productRepository repository.ProductRepository
	// orderRepository   repository.OrderRepository
	helper helper.Helper
}

func NewUserService(repo repository.UserRepository,
	cfg config.Config,
	// otp repository.OtpRepository,
	// inv repository.ProductRepository,
	// order repository.OrderRepository,
	h helper.Helper) UserService {

	return &userService{
		userRepo: repo,
		cfg:      cfg,
		// otpRepository:     otp,
		// productRepository: inv,
		// orderRepository:   order,
		helper: h,
	}
}

var InternalError = "Internal Server Error"
var ErrorHashingPassword = "Error In Hashing Password"

func (u *userService) Register(user models.UserDetails, ref string) (models.TokenUsers, error) {

	userExist := u.userRepo.CheckUserAvailability(user.Email, user.Phone)
	if userExist {
		return models.TokenUsers{}, models.ErrAlreadyExists
	}
	if user.Password != user.ConfirmPassword {
		return models.TokenUsers{}, models.ErrBadRequest
	}

	// referenceUser, err := u.userRepo.FindUserFromReference(ref)
	// if err != nil {
	// 	return models.TokenUsers{}, errors.New("cannot find reference user")
	// }

	// Hash password since details are validated

	hashedPassword, err := u.helper.PasswordHashing(user.Password)
	if err != nil {
		return models.TokenUsers{}, err
	}

	user.Password = hashedPassword

	referral, err := u.helper.GenerateRefferalCode()
	if err != nil {
		return models.TokenUsers{}, err
	}

	// add user details to the database
	userData, err := u.userRepo.Register(user, referral)
	if err != nil {
		return models.TokenUsers{}, err
	}
	user.Address.Name = user.Name
	user.Address.Phone = user.Phone
	user.Address.Default = true

	_, err = u.userRepo.AddAddress(userData.ID, user.Address)
	if err != nil {
		return models.TokenUsers{}, err
	}

	// crete a JWT token string for the user
	tokenString, err := u.helper.GenerateTokenClients(userData)
	if err != nil {
		return models.TokenUsers{}, err
	}

	// //credit 20 rupees to the user which is the source of the reference code
	// if err := u.userRepo.CreditReferencePointsToWallet(referenceUser); err != nil {
	// 	return models.TokenUsers{}, errors.New("error in crediting gift")
	// }

	//create new wallet for user
	// if _, err := u.orderRepository.CreateNewWallet(userData.Id); err != nil {
	// 	return models.TokenUsers{}, errors.New("errror in creating new wallet")
	// }

	return models.TokenUsers{
		Users: userData,
		Token: tokenString,
	}, nil
}

func (u *userService) Login(user models.UserLogin) (models.TokenUsers, error) {

	details, err := u.userRepo.FindUser(user)
	if err != nil {
		return models.TokenUsers{}, err
	}

	err = u.helper.CompareHashAndPassword(details.Password, user.Password)
	if err != nil {
		return models.TokenUsers{}, models.ErrInvalidPassword
	}

	userDetails := models.UserDetailsResponse{
		ID:        details.ID,
		Name:      details.Name,
		Email:     details.Email,
		Phone:     details.Phone,
		Username:  details.Username,
		Gender:    details.Gender,
		BirthDate: details.BirthDate,
	}

	tokenString, err := u.helper.GenerateTokenClients(userDetails)
	if err != nil {
		return models.TokenUsers{}, err
	}

	return models.TokenUsers{
		Users: userDetails,
		Token: tokenString,
	}, nil

}

func (i *userService) AddAddress(user_id uint, address models.Address) (models.Address, error) {

	addAddress, err := i.userRepo.AddAddress(user_id, address)
	if err != nil {
		return models.Address{}, err
	}

	return addAddress, nil

}

func (i *userService) UpdateAddress(user_id, address_id uint, address models.Address) (models.Address, error) {

	updateAddress, err := i.userRepo.UpdateAddress(user_id, address_id, address)
	if err != nil {
		return models.Address{}, err
	}

	return updateAddress, nil

}

func (i *userService) DeleteAddress(user_id, address_id uint) error {

	err := i.userRepo.DeleteAddress(user_id, address_id)
	if err != nil {
		return err
	}

	return nil

}

func (i *userService) GetAddresses(user_id uint) ([]models.Address, error) {

	addresses, err := i.userRepo.GetAddresses(user_id)
	if err != nil {
		return []models.Address{}, err
	}

	return addresses, nil

}

func (u *userService) GetUserDetails(id uint) (models.UserDetailsResponse, error) {

	details, err := u.userRepo.GetUserDetails(id)
	if err != nil {
		return models.UserDetailsResponse{}, err
	}

	return details, nil

}

func (u *userService) ChangePassword(id uint, old string, password string, repassword string) error {

	userPassword, err := u.userRepo.GetPassword(id)
	if err != nil {
		return errors.New(InternalError)
	}

	err = u.helper.CompareHashAndPassword(userPassword, old)
	if err != nil {
		return err
	}

	if password != repassword {
		return err
	}

	newpassword, err := u.helper.PasswordHashing(password)
	if err != nil {
		return err
	}

	return u.userRepo.ChangePassword(id, string(newpassword))

}

// func (u *userService) ForgotPasswordSend(phone string) error {

// 	ok := u.otpRepository.FindUserByMobileNumber(phone)
// 	if !ok {
// 		return errors.New("the user does not exist")
// 	}

// 	u.helper.TwilioSetup(u.cfg.ACCOUNTSID, u.cfg.AUTHTOKEN)
// 	_, err := u.helper.TwilioSendOTP(phone, u.cfg.SERVICESID)
// 	if err != nil {
// 		return errors.New("error ocurred while generating OTP")
// 	}

// 	return nil

// }

// func (u *USERUSECASE) ForgotPasswordVerifyAndChange(model models.ForgotVerify) error {
// 	u.helper.TwilioSetup(u.cfg.ACCOUNTSID, u.cfg.AUTHTOKEN)
// 	err := u.helper.TwilioVerifyOTP(u.cfg.SERVICESID, model.Otp, model.Phone)
// 	if err != nil {
// 		return errors.New("error while verifying")
// 	}

// 	id, err := u.userRepo.FindIdFromPhone(model.Phone)
// 	if err != nil {
// 		return errors.New("cannot find user from mobile number")
// 	}

// 	newpassword, err := u.helper.PasswordHashing(model.NewPassword)
// 	if err != nil {
// 		return errors.New("error in hashing password")
// 	}

// 	// if user is authenticated then change the password i the database
// 	if err := u.userRepo.ChangePassword(id, string(newpassword)); err != nil {
// 		return errors.New("could not change password")
// 	}

// 	return nil
// }

func (u *userService) EditProfile(user_id uint, profile models.EditProfile) (models.UserDetailsResponse, error) {

	result, err := u.userRepo.EditProfile(user_id, profile)
	if err != nil {
		return models.UserDetailsResponse{}, err
	}

	return result, nil

}

// func (u *USERUSECASE) GetMyReferenceLink(id uint) (string, error) {

// 	baseURL := "ahava.com/users/signup"

// 	referralCode, err := u.userRepo.GetReferralCodeFromID(id)
// 	if err != nil {
// 		return "", errors.New("error getting ref code")
// 	}

// 	referralLink := fmt.Sprintf("%s?ref=%s", baseURL, referralCode)

// 	//returning the link
// 	return referralLink, nil
// }
