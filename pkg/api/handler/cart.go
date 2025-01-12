package handler

import (
	services "ahava/pkg/service"
	"ahava/pkg/utils/models"
	"ahava/pkg/utils/response"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler interface {
	AddToCart(ctx *gin.Context)
	GetCart(ctx *gin.Context)
	RemoveFromCart(ctx *gin.Context)
	UpdateQuantityAdd(ctx *gin.Context)
	UpdateQuantityLess(ctx *gin.Context)
	UpdateQuantity(ctx *gin.Context)
	CheckOut(ctx *gin.Context)
}

type cartHandler struct {
	service services.CartService
}

func NewCartHandler(service services.CartService) CartHandler {
	return &cartHandler{
		service: service,
	}
}

func (i *cartHandler) AddToCart(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(uint)
	var model models.UpdateCartItem
	if err := ctx.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := i.service.AddToCart(user_id, model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the Cart", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added To cart", result, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *cartHandler) GetCart(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(uint)

	products, err := i.service.GetCart(user_id, []uint{})
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve cart", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all products in cart", products, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *cartHandler) RemoveFromCart(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(int)

	cart_id, err := strconv.Atoi(ctx.Param("cart_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := i.service.RemoveFromCart(uint(user_id), uint(cart_id)); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not remove from cart", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully Removed product from cart", nil, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *cartHandler) UpdateQuantityAdd(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(int)

	cart_id, err := strconv.Atoi(ctx.Param("cart_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	var model models.UpdateCartItem
	if err := ctx.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := i.service.UpdateQuantityAdd(uint(user_id), uint(cart_id), model.Quantity)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not Add the quantity", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added quantity", result, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *cartHandler) UpdateQuantityLess(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(int)

	cart_id, err := strconv.Atoi(ctx.Param("cart_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	var model models.UpdateCartItem
	if err := ctx.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := i.service.UpdateQuantityLess(uint(user_id), uint(cart_id), model.Quantity)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not subtract quantity", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully subtracted quantity", result, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *cartHandler) UpdateQuantity(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(int)

	cart_id, err := strconv.Atoi(ctx.Param("cart_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	var model models.UpdateCartItem
	if err := ctx.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	fmt.Println(model.Quantity)

	result, err := i.service.UpdateQuantity(uint(user_id), uint(cart_id), model.Quantity)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not update quantity", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated quantity", result, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *cartHandler) CheckOut(ctx *gin.Context) {
	user_id := ctx.MustGet("id").(int)

	var model models.CartCheckout
	if err := ctx.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	checkout, err := i.service.CheckOut(uint(user_id), model.CartIDs)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not open checkout", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", checkout, nil)
	ctx.JSON(http.StatusOK, successRes)
}
