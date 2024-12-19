package handler

import (
	services "ahava/pkg/usecase"
	"ahava/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	usecase services.WishlistUseCase
}

func NewWishlistHandler(use services.WishlistUseCase) *WishlistHandler {
	return &WishlistHandler{
		usecase: use,
	}
}

func (w *WishlistHandler) AddToWishlist(c *gin.Context) {

	user_id := c.MustGet("id").(int)

	product_id, err := strconv.Atoi(c.Query("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	result, err := w.usecase.AddToWishlist(user_id, product_id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add to Wishlist", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added To wishlist", result, nil)
	c.JSON(http.StatusOK, successRes)

}

func (w *WishlistHandler) RemoveFromWishlist(c *gin.Context) {

	user_id := c.MustGet("id").(int)

	product_id, err := strconv.Atoi(c.Query("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "check parameters properly", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	if err := w.usecase.RemoveFromWishlist(user_id, product_id); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not remove from wishlist", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully Removed product from wishlist", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func (w *WishlistHandler) GetWishList(c *gin.Context) {

	user_id := c.MustGet("id").(int)

	products, err := w.usecase.GetWishList(user_id)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", products, nil)
	c.JSON(http.StatusOK, successRes)
}
