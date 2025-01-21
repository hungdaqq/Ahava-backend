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
	ListAllOrders(ctx *gin.Context)
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
	// Get the user id from the context
	user_id := ctx.MustGet("id").(int)
	// Bind the request body to the model
	var orderDetails models.PlaceOrder
	if err := ctx.BindJSON(&orderDetails); err != nil {
		errorRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Perform place order operation
	orderDetails.UserID = uint(user_id)
	order, err := h.orderService.PlaceOrder(orderDetails)
	if err != nil {
		errorRes := response.ClientErrorResponse("Đặt hàng thất bại", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusCreated, "Đặt hàng thành công", order, nil)
	ctx.JSON(http.StatusCreated, successRes)
}

func (h *orderHandler) GetOrderDetails(ctx *gin.Context) {
	// Get the user id from the context
	user_id := ctx.MustGet("id").(int)
	// Get the order id from the query
	order_id, err := strconv.Atoi(ctx.Query("order_id"))
	if err != nil {
		errorRes := response.ClientErrorResponse("Request query problem", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Perform get order details operation
	order, err := h.orderService.GetOrderDetails(uint(user_id), uint(order_id))
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể lấy thông tin đơn hàng", nil, err)
		ctx.JSON(errorRes.StatusCode, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Lấy thông tin đơn hàng thành công", order, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (h *orderHandler) ListAllOrders(ctx *gin.Context) {
	// Get the limit and offset from the query
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		offset = 0
	}
	// Perform list orders operation
	orders, err := h.orderService.ListAllOrders(limit, offset)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể lấy danh sách đơn hàng", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Lấy danh sách đơn hàng thành công", orders, nil)
	ctx.JSON(http.StatusOK, successRes)
}
