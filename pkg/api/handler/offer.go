package handler

import (
	services "ahava/pkg/usecase"
	"ahava/pkg/utils/models"
	"ahava/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OfferHandler struct {
	usecase services.OfferUseCase
}

func NewOfferHandler(usecase services.OfferUseCase) *OfferHandler {
	return &OfferHandler{
		usecase: usecase,
	}
}

func (off *OfferHandler) AddNewOffer(c *gin.Context) {

	var model models.Offer
	if err := c.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := off.usecase.AddNewOffer(model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the Offer", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Successfully added Offer", result, nil)
	c.JSON(http.StatusOK, successRes)

}

func (off *OfferHandler) UpdateOffer(c *gin.Context) {

	product_id, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	var model models.Offer
	if err := c.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := off.usecase.UpdateOffer(uint(product_id), model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not update the Offer", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated Offer", result, nil)
	c.JSON(http.StatusOK, successRes)

}

func (o *OfferHandler) ExpireOffer(c *gin.Context) {

	product_id, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := o.usecase.ExpireOffer(uint(product_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Offer cannot be made expired", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully made offer expired", result, nil)
	c.JSON(http.StatusOK, successRes)
}

func (o *OfferHandler) GetOffers(c *gin.Context) {

	categories, err := o.usecase.GetOffers()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all offers", categories, nil)
	c.JSON(http.StatusOK, successRes)

}
