package handler

import (
	services "ahava/pkg/service"
	models "ahava/pkg/utils/models"
	"ahava/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WishlistHandler interface {
	AddToWishlist(ctx *gin.Context)
	RemoveFromWishlist(ctx *gin.Context)
	GetWishList(ctx *gin.Context)
}

type wishlistHandler struct {
	service services.WishlistService
}

func NewWishlistHandler(service services.WishlistService) WishlistHandler {
	return &wishlistHandler{
		service: service,
	}
}

func (h *wishlistHandler) AddToWishlist(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(int)

	var model models.AddToWishlist
	if err := ctx.BindJSON(&model); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := h.service.AddToWishlist(uint(user_id), model.ProductID)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add to Wishlist", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added To wishlist", result, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (h *wishlistHandler) RemoveFromWishlist(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(int)

	product_id, err := strconv.Atoi(ctx.Query("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	if err := h.service.RemoveFromWishlist(uint(user_id), uint(product_id)); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not remove from wishlist", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully Removed product from wishlist", nil, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (h *wishlistHandler) GetWishList(ctx *gin.Context) {

	user_id := ctx.MustGet("id").(int)

	products, err := h.service.GetWishList(uint(user_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", products, nil)
	ctx.JSON(http.StatusOK, successRes)
}
