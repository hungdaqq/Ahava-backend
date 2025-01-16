package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	services "ahava/pkg/service"
	"ahava/pkg/utils/models"
	"ahava/pkg/utils/response"
)

type UserHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
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

func (h *userHandler) Register(ctx *gin.Context) {
	// Bind the request body to the model
	var user models.UserDetails
	if err := ctx.BindJSON(&user); err != nil {
		errRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}
	// Validate the model
	err := validator.New().Struct(user)
	if err != nil {
		errRes := response.ClientErrorResponse("Constraints are not satisfied", nil, err)
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}
	// Get the reference from the query
	ref := ctx.Query("reference")
	// Perform register operation
	result, err := h.userService.Register(user, ref)
	if err != nil {
		errRes := response.ClientErrorResponse("Không thể đăng ký tài khoản", nil, err)
		ctx.JSON(http.StatusBadRequest, errRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusCreated, "Đăng ký tài khoản thành công", result, nil)
	ctx.JSON(http.StatusCreated, successRes)
}

func (h *userHandler) Login(ctx *gin.Context) {
	// Bind the request body to the model
	var user models.UserLogin
	if err := ctx.BindJSON(&user); err != nil {
		errRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}
	// Validate the model
	validator := validator.New()
	if err := validator.Struct(user); err != nil {
		errRes := response.ClientErrorResponse("Constraints are not satisfied", nil, err)
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}
	// Perform login operation
	result, err := h.userService.Login(user)
	if err != nil {
		errRes := response.ClientErrorResponse("Không thể đăng nhập", nil, err)
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Đăng nhập thành công", result, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *userHandler) AddAddress(ctx *gin.Context) {
	// Get the user id from the context
	user_id := ctx.MustGet("id").(int)
	// Bind the request body to the model
	var address models.Address
	if err := ctx.BindJSON(&address); err != nil {
		errRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}
	// Perform add operation
	result, err := i.userService.AddAddress(uint(user_id), address)
	if err != nil {
		errRes := response.ClientErrorResponse("Không thể thêm địa chỉ", nil, err)
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Thêm địa chỉ thành công", result, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *userHandler) UpdateAddress(ctx *gin.Context) {
	// Get the user id from the context
	user_id := ctx.MustGet("id").(int)
	// Get the address id from the url
	address_id, err := strconv.Atoi(ctx.Param("address_id"))
	if err != nil {
		errRes := response.ClientErrorResponse("Request parameter problem", nil, err)
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}
	// Bind the request body to the model
	var model models.Address
	if err := ctx.BindJSON(&model); err != nil {
		errRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}
	// Perform update operation
	result, err := i.userService.UpdateAddress(uint(user_id), uint(address_id), model)
	if err != nil {
		errRes := response.ClientErrorResponse("Không thể cập nhật địa chỉ", nil, err)
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Cập nhật địa chỉ thành công", result, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *userHandler) DeleteAddress(ctx *gin.Context) {
	// Get the user id from the context
	user_id := ctx.MustGet("id").(int)
	// Get the address id from the url
	address_id, err := strconv.Atoi(ctx.Param("address_id"))
	if err != nil {
		errRes := response.ClientErrorResponse("Request parameter problem", nil, err)
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}
	// Perform delete operation
	err = i.userService.DeleteAddress(uint(user_id), uint(address_id))
	if err != nil {
		errRes := response.ClientErrorResponse("Không thể xoá địa chỉ", nil, err)
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Xoá địa chỉ thành công", nil, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *userHandler) GetAddresses(ctx *gin.Context) {
	// Get the user id from the context
	user_id := ctx.MustGet("id").(int)
	// Perform get operation
	addresses, err := i.userService.GetAddresses(uint(user_id))
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể lấy danh sách địa chỉ", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Lấy danh sách địa chỉ thành công", addresses, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *userHandler) GetUserDetails(ctx *gin.Context) {
	// Get the user id from the context
	user_id := ctx.MustGet("id").(int)
	// Perform get operation
	details, err := i.userService.GetUserDetails(uint(user_id))
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể lấy thông tin người dùng", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Lấy thông tin người dùng thành công", details, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *userHandler) ChangePassword(ctx *gin.Context) {
	// Get the user id from the context
	user_id := ctx.MustGet("id").(int)
	// Bind the request body to the model
	var model models.ChangePassword
	if err := ctx.BindJSON(&model); err != nil {
		errRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}
	// Validate the model
	validator := validator.New()
	if err := validator.Struct(model); err != nil {
		errRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}
	// Perform change password operation
	if err := i.userService.ChangePassword(uint(user_id), model.Oldpassword, model.Password, model.Repassword); err != nil {
		errorRes := response.ClientErrorResponse("Không thể đổi mật khẩu", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Đổi mật khẩu thành công", nil, nil)
	ctx.JSON(http.StatusOK, successRes)
}

// func (i *userHandler) ForgotPasswordSend(ctx *gin.Context) {

// 	var model models.ForgotPasswordSend
// 	if err := ctx.BindJSON(&model); err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
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
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
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
	// Get the user id from the context
	user_id := ctx.MustGet("id").(int)
	// Bind the request body to the model
	var model models.EditProfile
	if err := ctx.BindJSON(&model); err != nil {
		errRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}
	// Perform edit operation
	result, err := i.userService.EditProfile(uint(user_id), model)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể sửa thông tin người dùng thành công", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Sửa thông tin người dùng thành công", result, nil)
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
