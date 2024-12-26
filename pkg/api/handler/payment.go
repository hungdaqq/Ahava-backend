package handler

import (
	"net/http"

	services "ahava/pkg/usecase"
	models "ahava/pkg/utils/models"
	response "ahava/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	orderUseCase services.PaymentUseCase
}

func NewPaymentHandler(usecase services.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		orderUseCase: usecase,
	}
}

func (h *PaymentHandler) CreateQR(c *gin.Context) {

	var model models.CreateQR
	err := c.BindJSON(&model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "parameter problem", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := h.orderUseCase.CreateQR(model.Amount)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not create QR", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully placed the order", result, nil)
	c.JSON(http.StatusOK, successRes)
}
