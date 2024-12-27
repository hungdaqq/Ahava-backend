package usecase

import (
	"errors"

	"ahava/pkg/config"
	helper "ahava/pkg/helper"
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
)

type UserUseCase interface {
	UserSignUp(user models.UserDetails, ref string) (models.TokenUsers, error)
	LoginHandler(user models.UserLogin) (models.TokenUsers, error)
	AddAddress(user_id uint, address models.Address) (models.Address, error)
	GetAddresses(user_id uint) ([]models.Address, error)
	UpdateAddress(address_id uint, address models.Address) (models.Address, error)
	DeleteAddress(address_id uint) error

	GetUserDetails(user_id uint) (models.UserDetailsResponse, error)

	ChangePassword(user_id uint, old string, password string, repassword string) error
	// ForgotPasswordSend(phone string) error
	// ForgotPasswordVerifyAndChange(model models.ForgotVerify) error

	EditProfile(user_id uint, name, email, phone string) (models.UserDetailsResponse, error)

	// GetMyReferenceLink(id uint) (string, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
	cfg      config.Config
	// otpRepository     repository.OtpRepository
	// productRepository repository.ProductRepository
	// orderRepository   repository.OrderRepository
	helper helper.Helper
}

func NewUserUseCase(repo repository.UserRepository,
	cfg config.Config,
	// otp repository.OtpRepository,
	// inv repository.ProductRepository,
	// order repository.OrderRepository,
	h helper.Helper) *userUseCase {

	return &userUseCase{
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

func (u *userUseCase) UserSignUp(user models.UserDetails, ref string) (models.TokenUsers, error) {
	// Check whether the user already exist. If yes, show the error message, since this is signUp
	userExist := u.userRepo.CheckUserAvailability(user.Email, user.Phone)
	if userExist {
		return models.TokenUsers{}, errors.New("user already exist, sign in")
	}
	if user.Password != user.ConfirmPassword {
		return models.TokenUsers{}, errors.New("password does not match")
	}

	// referenceUser, err := u.userRepo.FindUserFromReference(ref)
	// if err != nil {
	// 	return models.TokenUsers{}, errors.New("cannot find reference user")
	// }

	// Hash password since details are validated

	hashedPassword, err := u.helper.PasswordHashing(user.Password)
	if err != nil {
		return models.TokenUsers{}, errors.New(ErrorHashingPassword)
	}

	user.Password = hashedPassword

	referral, err := u.helper.GenerateRefferalCode()
	if err != nil {
		return models.TokenUsers{}, errors.New(InternalError)
	}

	// add user details to the database
	userData, err := u.userRepo.UserSignUp(user, referral)
	if err != nil {
		return models.TokenUsers{}, errors.New("could not add the user")
	}
	user.Address.Name = user.Name
	user.Address.Phone = user.Phone
	user.Address.Default = true

	_, err = u.userRepo.AddAddress(userData.ID, user.Address)
	if err != nil {
		return models.TokenUsers{}, errors.New("could not add the address")
	}

	// crete a JWT token string for the user
	tokenString, err := u.helper.GenerateTokenClients(userData)
	if err != nil {
		return models.TokenUsers{}, errors.New("could not create token due to some internal error")
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

func (u *userUseCase) LoginHandler(user models.UserLogin) (models.TokenUsers, error) {

	// checking if a username exist with this email address
	ok := u.userRepo.CheckUserAvailability(user.Email, user.Phone)
	if !ok {
		return models.TokenUsers{}, errors.New("the user does not exist")
	}

	isBlocked, err := u.userRepo.UserBlockStatus(user.Email)
	if err != nil {
		return models.TokenUsers{}, errors.New(InternalError)
	}

	if isBlocked {
		return models.TokenUsers{}, errors.New("user is blocked by admin")
	}

	// Get the user details in order to check the password, in this case ( The same function can be reused in future )
	user_details, err := u.userRepo.FindUserByEmail(user)
	if err != nil {
		return models.TokenUsers{}, errors.New(InternalError)
	}

	err = u.helper.CompareHashAndPassword(user_details.Password, user.Password)
	if err != nil {
		return models.TokenUsers{}, errors.New("password incorrect")
	}

	var userDetails models.UserDetailsResponse

	userDetails.ID = user_details.ID
	userDetails.Name = user_details.Name
	userDetails.Email = user_details.Email
	userDetails.Phone = user_details.Phone

	tokenString, err := u.helper.GenerateTokenClients(userDetails)
	if err != nil {
		return models.TokenUsers{}, errors.New("could not create token")
	}

	return models.TokenUsers{
		Users: userDetails,
		Token: tokenString,
	}, nil

}

func (i *userUseCase) AddAddress(user_id uint, address models.Address) (models.Address, error) {

	addAddress, err := i.userRepo.AddAddress(user_id, address)
	if err != nil {
		return models.Address{}, errors.New("error in adding address")
	}

	return addAddress, nil

}

func (i *userUseCase) UpdateAddress(address_id uint, address models.Address) (models.Address, error) {

	updateAddress, err := i.userRepo.UpdateAddress(address_id, address)
	if err != nil {
		return models.Address{}, errors.New("error in updating address")
	}

	return updateAddress, nil

}

func (i *userUseCase) DeleteAddress(address_id uint) error {

	err := i.userRepo.DeleteAddress(address_id)
	if err != nil {
		return errors.New("error in deleting address")
	}

	return nil

}

func (i *userUseCase) GetAddresses(user_id uint) ([]models.Address, error) {

	addresses, err := i.userRepo.GetAddresses(user_id)
	if err != nil {
		return []models.Address{}, errors.New("error in getting addresses")
	}

	return addresses, nil

}

func (u *userUseCase) GetUserDetails(id uint) (models.UserDetailsResponse, error) {

	details, err := u.userRepo.GetUserDetails(id)
	if err != nil {
		return models.UserDetailsResponse{}, errors.New("error in getting details")
	}

	return details, nil

}

func (u *userUseCase) ChangePassword(id uint, old string, password string, repassword string) error {

	userPassword, err := u.userRepo.GetPassword(id)
	if err != nil {
		return errors.New(InternalError)
	}

	err = u.helper.CompareHashAndPassword(userPassword, old)
	if err != nil {
		return errors.New("password incorrect")
	}

	if password != repassword {
		return errors.New("passwords does not match")
	}

	newpassword, err := u.helper.PasswordHashing(password)
	if err != nil {
		return errors.New("error in hashing password")
	}

	return u.userRepo.ChangePassword(id, string(newpassword))

}

// func (u *userUseCase) ForgotPasswordSend(phone string) error {

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

// func (u *userUseCase) ForgotPasswordVerifyAndChange(model models.ForgotVerify) error {
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

func (u *userUseCase) EditProfile(user_id uint, name, email, phone string) (models.UserDetailsResponse, error) {

	result, err := u.userRepo.EditProfile(user_id, name, email, phone)
	if err != nil {
		return models.UserDetailsResponse{}, errors.New("could not change")
	}

	return result, nil

}

// func (u *userUseCase) GetMyReferenceLink(id uint) (string, error) {

// 	baseURL := "ahava.com/users/signup"

// 	referralCode, err := u.userRepo.GetReferralCodeFromID(id)
// 	if err != nil {
// 		return "", errors.New("error getting ref code")
// 	}

// 	referralLink := fmt.Sprintf("%s?ref=%s", baseURL, referralCode)

// 	//returning the link
// 	return referralLink, nil
// }
