package handler

import (
	services "ahava/pkg/service"
	"ahava/pkg/utils/models"
	"ahava/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OfferHandler interface {
	AddNewOffer(ctx *gin.Context)
	UpdateOffer(ctx *gin.Context)
	ExpireOffer(ctx *gin.Context)
	GetOffers(ctx *gin.Context)
}

type offerHandler struct {
	service services.OfferService
}

func NewOfferHandler(service services.OfferService) OfferHandler {
	return &offerHandler{
		service: service,
	}
}

func (h *offerHandler) AddNewOffer(ctx *gin.Context) {

	var model models.Offer
	if err := ctx.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := h.service.AddNewOffer(model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the Offer", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Successfully added Offer", result, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (h *offerHandler) UpdateOffer(ctx *gin.Context) {

	product_id, err := strconv.Atoi(ctx.Param("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	var model models.Offer
	if err := ctx.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := h.service.UpdateOffer(uint(product_id), model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not update the Offer", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated Offer", result, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (o *offerHandler) ExpireOffer(ctx *gin.Context) {

	product_id, err := strconv.Atoi(ctx.Param("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := o.service.ExpireOffer(uint(product_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Offer cannot be made expired", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully made offer expired", result, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (o *offerHandler) GetOffers(ctx *gin.Context) {

	categories, err := o.service.GetOffers()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all offers", categories, nil)
	ctx.JSON(http.StatusOK, successRes)

}
