package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	services "ahava/pkg/usecase"
	"ahava/pkg/utils/models"
	"ahava/pkg/utils/response"
)

type UserHandler struct {
	userUseCase services.UserUseCase
}

func NewUserHandler(usecase services.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
	}
}

func (u *UserHandler) UserSignUp(c *gin.Context) {

	var user models.UserDetails
	// bind the user details to the struct
	if err := c.BindJSON(&user); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// checking whether the data sent by the user has all the correct constraints specified by Users struct
	err := validator.New().Struct(user)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest,
			errRes)
		return
	}

	//if the user wants to mention the referral code of other user
	ref := c.Query("reference")

	// business logic goes inside this function
	userCreated, err := u.userUseCase.UserSignUp(user, ref)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "User could not signed up", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "User successfully signed up", userCreated, nil)
	c.JSON(http.StatusCreated, successRes)

}

func (u *UserHandler) LoginHandler(c *gin.Context) {

	var user models.UserLogin

	if err := c.BindJSON(&user); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := validator.New().Struct(user)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	user_details, err := u.userUseCase.LoginHandler(user)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "User could not be logged in", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "User successfully logged in", user_details, nil)
	c.JSON(http.StatusOK, successRes)

}

func (i *UserHandler) AddAddress(c *gin.Context) {

	user_id := c.MustGet("id").(int)
	var address models.Address
	if err := c.BindJSON(&address); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := i.userUseCase.AddAddress(user_id, address)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the address", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added address", result, nil)
	c.JSON(http.StatusOK, successRes)

}

func (i *UserHandler) GetAddresses(c *gin.Context) {

	user_id := c.MustGet("id").(int)
	addresses, err := i.userUseCase.GetAddresses(user_id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", addresses, nil)
	c.JSON(http.StatusOK, successRes)
}

func (i *UserHandler) GetUserDetails(c *gin.Context) {

	user_id := c.MustGet("id").(int)
	details, err := i.userUseCase.GetUserDetails(user_id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", details, nil)
	c.JSON(http.StatusOK, successRes)
}

func (i *UserHandler) ChangePassword(c *gin.Context) {

	user_id := c.MustGet("id").(int)
	var ChangePassword models.ChangePassword
	if err := c.BindJSON(&ChangePassword); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := i.userUseCase.ChangePassword(user_id, ChangePassword.Oldpassword, ChangePassword.Password, ChangePassword.Repassword); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not change the password", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "password changed Successfully ", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

// func (i *UserHandler) ForgotPasswordSend(c *gin.Context) {

// 	var model models.ForgotPasswordSend
// 	if err := c.BindJSON(&model); err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}
// 	err := i.userUseCase.ForgotPasswordSend(model.Phone)
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not send OTP", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	successRes := response.ClientResponse(http.StatusOK, "OTP sent successfully", nil, nil)
// 	c.JSON(http.StatusOK, successRes)

// }

// func (i *UserHandler) ForgotPasswordVerifyAndChange(c *gin.Context) {

// 	var model models.ForgotVerify
// 	if err := c.BindJSON(&model); err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	err := i.userUseCase.ForgotPasswordVerifyAndChange(model)
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not verify OTP", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	successRes := response.ClientResponse(http.StatusOK, "Successfully Changed the password", nil, nil)
// 	c.JSON(http.StatusOK, successRes)

// }

func (i *UserHandler) EditProfile(c *gin.Context) {

	user_id := c.MustGet("id").(int)
	var model models.EditProfile
	if err := c.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := i.userUseCase.EditProfile(user_id, model.Name, model.Email, model.Phone)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not edit profile", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully edited profile", result, nil)
	c.JSON(http.StatusOK, successRes)

}

func (i *UserHandler) GetMyReferenceLink(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	link, err := i.userUseCase.GetMyReferenceLink(id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve referral link", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all products in cart", link, nil)
	c.JSON(http.StatusOK, successRes)
}
