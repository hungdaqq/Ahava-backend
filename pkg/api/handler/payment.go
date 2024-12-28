package handler

import (
	"net/http"

	services "ahava/pkg/usecase"
	models "ahava/pkg/utils/models"
	response "ahava/pkg/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

	user_id := c.MustGet("id").(int)

	var model models.CreateQR
	err := c.BindJSON(&model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "parameter problem", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = validator.New().Struct(model)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	result, err := h.orderUseCase.CreateSePayQR(uint(user_id), model.OrderID, model.Amount)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not create QR", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully placed the order", result, nil)
	c.JSON(http.StatusOK, successRes)
}

func (h *PaymentHandler) Webhook(c *gin.Context) {

	var model models.Transaction
	err := c.BindJSON(&model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "parameter problem", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = h.orderUseCase.Webhook(model)
	if err != nil {
		errorRes := response.ClientWebhookResponse(false)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientWebhookResponse(true)
	c.JSON(http.StatusCreated, successRes)
}
