package handler

import (
	services "ahava/pkg/usecase"
	"ahava/pkg/utils/models"
	"ahava/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	usecase services.CartUseCase
}

func NewCartHandler(usecase services.CartUseCase) *CartHandler {
	return &CartHandler{
		usecase: usecase,
	}
}

func (i *CartHandler) AddToCart(c *gin.Context) {

	user_id := c.MustGet("id").(int)
	var model models.UpdateCartItem
	if err := c.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	result, err := i.usecase.AddToCart(uint(user_id), uint(model.ProductID), model.Quantity)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the Cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added To cart", result, nil)
	c.JSON(http.StatusOK, successRes)
}

func (i *CartHandler) GetCart(c *gin.Context) {

	user_id := c.MustGet("id").(int)

	products, err := i.usecase.GetCart(uint(user_id), []uint{})
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all products in cart", products, nil)
	c.JSON(http.StatusOK, successRes)
}

func (i *CartHandler) RemoveFromCart(c *gin.Context) {

	cart_id, err := strconv.Atoi(c.Param("cart_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	if err := i.usecase.RemoveFromCart(uint(cart_id)); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not remove from cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully Removed product from cart", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (i *CartHandler) UpdateQuantityAdd(c *gin.Context) {

	cart_id, err := strconv.Atoi(c.Param("cart_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	var model models.UpdateCartItem
	if err := c.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := i.usecase.UpdateQuantityAdd(uint(cart_id), model.Quantity)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not Add the quantity", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added quantity", result, nil)
	c.JSON(http.StatusOK, successRes)
}

func (i *CartHandler) UpdateQuantityLess(c *gin.Context) {

	cart_id, err := strconv.Atoi(c.Param("cart_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	var model models.UpdateCartItem
	if err := c.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := i.usecase.UpdateQuantityLess(uint(cart_id), model.Quantity)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not subtract quantity", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully subtracted quantity", result, nil)
	c.JSON(http.StatusOK, successRes)
}

func (i *CartHandler) UpdateQuantity(c *gin.Context) {

	cart_id, err := strconv.Atoi(c.Param("cart_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	var model models.UpdateCartItem
	if err := c.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := i.usecase.UpdateQuantity(uint(cart_id), model.Quantity)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not  subtract quantity", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully subtracted quantity", result, nil)
	c.JSON(http.StatusOK, successRes)
}

func (i *CartHandler) CheckOut(c *gin.Context) {
	user_id := c.MustGet("id").(int)

	var model models.CartCheckout
	if err := c.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	checkout, err := i.usecase.CheckOut(uint(user_id), model.CartIDs)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not open checkout", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", checkout, nil)
	c.JSON(http.StatusOK, successRes)
}
