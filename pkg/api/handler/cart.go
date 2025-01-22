package handler

import (
	services "ahava/pkg/service"
	"ahava/pkg/utils/models"
	"ahava/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CartHandler interface {
	AddToCart(ctx *gin.Context)
	GetCart(ctx *gin.Context)
	RemoveFromCart(ctx *gin.Context)
	UpdateQuantity(ctx *gin.Context)
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
	// Get the user id from the context
	user_id := ctx.MustGet("id").(int)
	// Bind the request body to the model
	var model models.UpdateCartItem
	if err := ctx.BindJSON(&model); err != nil {
		errorRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Validate the model
	validator := validator.New()
	if err := validator.Struct(model); err != nil {
		errRes := response.ClientErrorResponse("Constraints not satisfied", nil, err)
		ctx.JSON(http.StatusBadRequest, errRes)
		return
	}
	// Perform add to cart operation
	result, err := i.service.AddToCart(uint(user_id), model)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể thêm sản phẩm vào giỏ hàng", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusCreated, "Thêm sản phẩm vào giỏ hàng thành công", result, nil)
	ctx.JSON(http.StatusCreated, successRes)
}

func (i *cartHandler) GetCart(ctx *gin.Context) {
	// Get the user id from the context
	user_id := ctx.MustGet("id").(int)
	// Perform get cart operation
	products, err := i.service.GetCart(uint(user_id), []uint{})
	if err != nil {
		errorRes := response.ClientErrorResponse("Lấy giỏ hàng thành công", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Không thể lấy giỏ hàng", products, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *cartHandler) RemoveFromCart(ctx *gin.Context) {
	// Get the user id from the context
	user_id := ctx.MustGet("id").(int)
	// Get the cart id from the params
	cart_id, err := strconv.Atoi(ctx.Param("cart_id"))
	if err != nil {
		errorRes := response.ClientErrorResponse("Request parameter problem", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Perform remove from cart operation
	if err := i.service.RemoveFromCart(uint(user_id), uint(cart_id)); err != nil {
		errorRes := response.ClientErrorResponse("Không thể xoá sản phẩm khỏi giỏ hàng", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Xoá sản phẩm khỏi giỏ hàng thành công", nil, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (i *cartHandler) UpdateQuantity(ctx *gin.Context) {
	// Get the user id from the context
	user_id := ctx.MustGet("id").(int)
	// Get the cart id from the params
	cart_id, err := strconv.Atoi(ctx.Param("cart_id"))
	if err != nil {
		errorRes := response.ClientErrorResponse("Request parameter problem", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Bind the request body to the model
	var model models.UpdateCartItem
	if err := ctx.BindJSON(&model); err != nil {
		errorRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Validate the model
	result, err := i.service.UpdateQuantity(uint(user_id), uint(cart_id), model.Quantity)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể cập nhật số lượng", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Cập nhật số lượng thành công", result, nil)
	ctx.JSON(http.StatusOK, successRes)
}

