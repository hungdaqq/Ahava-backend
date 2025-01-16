package handler

import (
	"net/http"
	"strconv"
	"time"

	"ahava/pkg/helper"
	services "ahava/pkg/service"
	models "ahava/pkg/utils/models"

	response "ahava/pkg/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AdminHandler interface {
	Login(ctx *gin.Context)
	BlockUser(ctx *gin.Context)
	UnBlockUser(ctx *gin.Context)
	GetUsers(ctx *gin.Context)
	NewPaymentMethod(ctx *gin.Context)
	ListPaymentMethods(ctx *gin.Context)
	DeletePaymentMethod(ctx *gin.Context)
	ValidateRefreshTokenAndCreateNewAccess(ctx *gin.Context)
}

type adminHandler struct {
	adminService services.AdminService
}

func NewAdminHandler(service services.AdminService) AdminHandler {
	return &adminHandler{
		adminService: service,
	}
}

func (ad *adminHandler) Login(ctx *gin.Context) {
	// Bind the request body to the model
	var adminDetails models.AdminLogin
	if err := ctx.BindJSON(&adminDetails); err != nil {
		errorRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Perform login operation
	admin, err := ad.adminService.Login(adminDetails)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể đăng nhập admin", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	ctx.Set("Access", admin.AccessToken)
	ctx.Set("Refresh", admin.RefreshToken)
	successRes := response.ClientResponse(http.StatusOK, "Đặng nhập admin thành công", admin, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (ad *adminHandler) BlockUser(ctx *gin.Context) {
	// Get the user id from the context
	user_id, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		errorRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Perform block user operation
	err = ad.adminService.BlockUser(uint(user_id))
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể chặn người dùng", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Chặn người dùng thành công", nil, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (ad *adminHandler) UnBlockUser(ctx *gin.Context) {
	// Get the user id from the context
	user_id, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		errorRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Perform unblock user operation
	err = ad.adminService.UnBlockUser(uint(user_id))
	if err != nil {
		errorRes := response.ClientErrorResponse("Không bỏ thể chặn người dùng", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Bỏ chặn người dùng thành công", nil, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (ad *adminHandler) GetUsers(ctx *gin.Context) {

	pageStr := ctx.Query("page")
	page, err := strconv.Atoi(pageStr)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	users, err := ad.adminService.GetUsers(page)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully retrieved the users", users, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (i *adminHandler) NewPaymentMethod(ctx *gin.Context) {

	var method models.NewPaymentMethod
	if err := ctx.BindJSON(&method); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err := i.adminService.NewPaymentMethod(method.PaymentMethod)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the payment method", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added Payment Method", nil, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (a *adminHandler) ListPaymentMethods(ctx *gin.Context) {

	categories, err := a.adminService.ListPaymentMethods()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all payment methods", categories, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (a *adminHandler) DeletePaymentMethod(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = a.adminService.DeletePaymentMethod(uint(id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "error in deleting data", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully deleted the Category", nil, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (a *adminHandler) ValidateRefreshTokenAndCreateNewAccess(ctx *gin.Context) {

	refreshToken := ctx.Request.Header.Get("RefreshToken")

	// Check if the refresh token is valid.
	_, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("refreshsecret"), nil
	})
	if err != nil {
		// The refresh token is invalid.
		ctx.AbortWithError(401, models.ErrInvalidToken)
		return
	}

	claims := &helper.AuthCustomClaims{
		Role: "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newAccessToken, err := token.SignedString([]byte("accesssecret"))
	if err != nil {
		ctx.AbortWithError(500, models.ErrCreateToken)
	}

	ctx.JSON(200, newAccessToken)
}
