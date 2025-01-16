package handler

import (
	"net/http"

	services "ahava/pkg/service"
	models "ahava/pkg/utils/models"
	response "ahava/pkg/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PaymentHandler interface {
	CreateQR(ctx *gin.Context)
	Webhook(ctx *gin.Context)
}

type paymentHandler struct {
	paymentService services.PaymentService
}

func NewPaymentHandler(service services.PaymentService) PaymentHandler {
	return &paymentHandler{
		paymentService: service,
	}
}

func (h *paymentHandler) CreateQR(c *gin.Context) {
	// Get the user id from the context
	user_id := c.MustGet("id").(int)
	// Bind the request body to the model
	var model models.CreateQR
	err := c.BindJSON(&model)
	if err != nil {
		errorRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Validate the model
	err = validator.New().Struct(model)
	if err != nil {
		errRes := response.ClientErrorResponse("Constraints not satisfied", nil, err)
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	// Perform create QR operation
	result, err := h.paymentService.CreateSePayQR(uint(user_id), model.OrderID, model.Amount)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể tạo QR", nil, err)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusCreated, "Tạo QR thành công", result, nil)
	c.JSON(http.StatusCreated, successRes)
}

func (h *paymentHandler) Webhook(c *gin.Context) {
	// Bind the request body to the model
	var model models.Transaction
	err := c.BindJSON(&model)
	if err != nil {
		errorRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Perform webhook operation
	err = h.paymentService.Webhook(model)
	if err != nil {
		errorRes := response.ClientWebhookResponse(false)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientWebhookResponse(true)
	c.JSON(http.StatusCreated, successRes)
}
