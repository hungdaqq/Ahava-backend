package handler

import (
	"net/http"
	"strconv"

	services "ahava/pkg/usecase"
	models "ahava/pkg/utils/models"
	response "ahava/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUseCase services.OrderUseCase
}

func NewOrderHandler(usecase services.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: usecase,
	}
}

func (h *OrderHandler) PlaceOrder(c *gin.Context) {

	user_id := c.MustGet("id").(int)

	var orderDetails models.PlaceOrder
	if err := c.BindJSON(&orderDetails); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	orderDetails.UserID = uint(user_id)

	order, err := h.orderUseCase.PlaceOrder(orderDetails)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not place the order", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully placed the order", order, nil)
	c.JSON(http.StatusOK, successRes)
}

func (h *OrderHandler) GetOrderPaidStatus(c *gin.Context) {

	order_id, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "parameter problem", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	status, err := h.orderUseCase.GetOrderPaidStatus(uint(order_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get the order status", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully fetched the order status", status, nil)
	c.JSON(http.StatusOK, successRes)
}

func (h *OrderHandler) GetOrderDetails(c *gin.Context) {

	order_id, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "parameter problem", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	order, err := h.orderUseCase.GetOrderDetails(uint(order_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not get the order details", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully fetched the order details", order, nil)
	c.JSON(http.StatusOK, successRes)
}
