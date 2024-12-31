package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	services "ahava/pkg/service"
	"ahava/pkg/utils/models"
	"ahava/pkg/utils/response"
)

type UserHandler interface {
	UserSignUp(ctx *gin.Context)
	LoginHandler(ctx *gin.Context)
	AddAddress(ctx *gin.Context)
	UpdateAddress(ctx *gin.Context)
	DeleteAddress(ctx *gin.Context)
	GetAddresses(ctx *gin.Context)
	GetUserDetails(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	// ForgotPasswordSend(ctx *gin.Context)
	// ForgotPasswordVerifyAndChange(ctx *gin.Context)
	EditProfile(ctx *gin.Context)
	// GetMyReferenceLink(ctx *gin.Context)
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(service services.UserService) UserHandler {
	return &userHandler{
		userService: service,
	}
}

func (h *userHandler) UserSignUp(ctx *gin.Context) {

	var user models.UserDetails
	if err := ctx.BindJSON(&user); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := validator.New().Struct(user)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		ctx.JSON(http.StatusBadRequest,
			errRes)
		return
	}

	ref := ctx.Query("reference")

	userCreated, err := h.userService.UserSignUp(user, ref)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "User could not signed up", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "User successfully signed up", userCreated, nil)
	ctx.JSON(http.StatusCreated, successRes)

}

func (h *userHandler) LoginHandler(ctx *gin.Context) {

	var user models.UserLogin

	if err := ctx.BindJSON(&user); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errRes)
		return
	}

	user_details, err := h.userService.LoginHandler(user)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "User could not be logged in", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "User successfully logged in", user_details, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (i *userHandler) AddAddress(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(int)
	var address models.Address
	if err := ctx.BindJSON(&address); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := i.userService.AddAddress(uint(user_id), address)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the address", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added address", result, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (i *userHandler) UpdateAddress(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(int)

	address_id, err := strconv.Atoi(ctx.Param("address_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "parameter problem", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	var model models.Address
	if err := ctx.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := i.userService.UpdateAddress(uint(user_id), uint(address_id), model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not update the address", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated address", result, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *userHandler) DeleteAddress(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(int)

	address_id, err := strconv.Atoi(ctx.Param("address_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "parameter problem", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = i.userService.DeleteAddress(uint(user_id), uint(address_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not delete the address", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully deleted address", nil, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *userHandler) GetAddresses(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(int)
	addresses, err := i.userService.GetAddresses(uint(user_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", addresses, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *userHandler) GetUserDetails(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(int)
	details, err := i.userService.GetUserDetails(uint(user_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", details, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *userHandler) ChangePassword(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(int)
	var ChangePassword models.ChangePassword
	if err := ctx.BindJSON(&ChangePassword); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := i.userService.ChangePassword(uint(user_id), ChangePassword.Oldpassword, ChangePassword.Password, ChangePassword.Repassword); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not change the password", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "password changed Successfully ", nil, nil)
	ctx.JSON(http.StatusOK, successRes)

}

// func (i *userHandler) ForgotPasswordSend(ctx *gin.Context) {

// 	var model models.ForgotPasswordSend
// 	if err := ctx.BindJSON(&model); err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
// 		ctx.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}
// 	err := i.userService.ForgotPasswordSend(model.Phone)
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not send OTP", nil, err.Error())
// 		ctx.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	successRes := response.ClientResponse(http.StatusOK, "OTP sent successfully", nil, nil)
// 	ctx.JSON(http.StatusOK, successRes)

// }

// func (i *userHandler) ForgotPasswordVerifyAndChange(ctx *gin.Context) {

// 	var model models.ForgotVerify
// 	if err := ctx.BindJSON(&model); err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
// 		ctx.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	err := i.userService.ForgotPasswordVerifyAndChange(model)
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not verify OTP", nil, err.Error())
// 		ctx.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	successRes := response.ClientResponse(http.StatusOK, "Successfully Changed the password", nil, nil)
// 	ctx.JSON(http.StatusOK, successRes)

// }

func (i *userHandler) EditProfile(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(int)
	var model models.EditProfile
	if err := ctx.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := i.userService.EditProfile(uint(user_id), model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not edit profile", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully edited profile", result, nil)
	ctx.JSON(http.StatusOK, successRes)

}

// func (i *userHandler) GetMyReferenceLink(ctx *gin.Context) {
// 	id, err := strconv.Atoi(ctx.Query("id"))
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
// 		ctx.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	link, err := i.userService.GetMyReferenceLink(id)
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve referral link", nil, err.Error())
// 		ctx.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}
// 	successRes := response.ClientResponse(http.StatusOK, "Successfully got all products in cart", link, nil)
// 	ctx.JSON(http.StatusOK, successRes)
// }

// FetchProvinces fetches the list of provinces from the API.
func FetchProvinceAPI(baseURL string, model interface{}) (interface{}, error) {

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	err = json.Unmarshal(body, &model)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return model, nil
}
