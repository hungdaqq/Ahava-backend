package handler

import (
	"net/http"
	"strconv"

	services "ahava/pkg/service"
	models "ahava/pkg/utils/models"
	response "ahava/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

type OrderHandler interface {
	PlaceOrder(ctx *gin.Context)
	GetOrderDetails(ctx *gin.Context)
}

type orderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(service services.OrderService) OrderHandler {
	return &orderHandler{
		orderService: service,
	}
}

func (h *orderHandler) PlaceOrder(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(int)

	var orderDetails models.PlaceOrder
	if err := ctx.BindJSON(&orderDetails); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	orderDetails.UserID = uint(user_id)

	order, err := h.orderService.PlaceOrder(orderDetails)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not place the order", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully placed the order", order, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (h *orderHandler) GetOrderDetails(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(int)

	order_id, err := strconv.Atoi(ctx.Query("order_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "parameter problem", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	order, err := h.orderService.GetOrderDetails(uint(user_id), uint(order_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get the order details", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully fetched the order details", order, nil)
	ctx.JSON(http.StatusOK, successRes)
}
