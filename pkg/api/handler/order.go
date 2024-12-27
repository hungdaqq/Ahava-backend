package handler

import (
	"net/http"

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
