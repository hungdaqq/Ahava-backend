package handler

import (
	"net/http"
	"strconv"
	"time"

	"ahava/pkg/helper"
	services "ahava/pkg/service"
	errors "ahava/pkg/utils/errors"
	models "ahava/pkg/utils/models"

	response "ahava/pkg/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AdminHandler interface {
	LoginHandler(ctx *gin.Context)
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

func (ad *adminHandler) LoginHandler(ctx *gin.Context) { // login handler for the admin

	// var adminDetails models.AdminLogin
	var adminDetails models.AdminLogin
	if err := ctx.BindJSON(&adminDetails); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "details not in correct format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errRes)
		return
	}

	admin, err := ad.adminService.LoginHandler(adminDetails)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "cannot authenticate user", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errRes)
		return
	}

	ctx.Set("Access", admin.AccessToken)
	ctx.Set("Refresh", admin.RefreshToken)

	successRes := response.ClientResponse(http.StatusOK, "Admin authenticated successfully", admin, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (ad *adminHandler) BlockUser(ctx *gin.Context) {

	user_id, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "user_id not in right format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = ad.adminService.BlockUser(uint(user_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "user could not be blocked", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully blocked the user", nil, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (ad *adminHandler) UnBlockUser(ctx *gin.Context) {

	user_id, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "user_id not in right format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = ad.adminService.UnBlockUser(uint(user_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "user could not be unblocked", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully unblocked the user", nil, nil)
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
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
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
		errorRes := response.ClientResponse(http.StatusInternalServerError, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all payment methods", categories, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (a *adminHandler) DeletePaymentMethod(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
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
		ctx.AbortWithError(401, errors.ErrInvalidToken)
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
		ctx.AbortWithError(500, errors.ErrCreateToken)
	}

	ctx.JSON(200, newAccessToken)
}
