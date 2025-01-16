package handler

import (
	services "ahava/pkg/service"
	models "ahava/pkg/utils/models"
	"ahava/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	// Get the user id from the context
	user_id := ctx.MustGet("id").(int)
	// Bind the request body to the model
	var model models.AddToWishlist
	if err := ctx.BindJSON(&model); err != nil {
		errorRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Validate the model
	validator := validator.New()
	if err := validator.Struct(model); err != nil {
		errorRes := response.ClientErrorResponse("Constraints are not satisfied", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Perform add to wishlist operation
	result, err := h.service.AddToWishlist(uint(user_id), model)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể thêm sản phẩm vào yêu thích", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Thêm sản phẩm vào yêu thích thành công", result, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (h *wishlistHandler) RemoveFromWishlist(ctx *gin.Context) {
	// Get the user id from the context
	user_id := ctx.MustGet("id").(int)
	// Get the wishlist id from the context
	wishlist_id, err := strconv.Atoi(ctx.Param("wishlist_id"))
	if err != nil {
		errRes := response.ClientErrorResponse("Request parameter problem", nil, err)
		ctx.JSON(errRes.StatusCode, errRes)
		return
	}
	// Perform remove from wishlist operation
	if err := h.service.RemoveFromWishlist(uint(user_id), uint(wishlist_id)); err != nil {
		errorRes := response.ClientErrorResponse("Không thể bỏ yêu thích", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Bỏ yêu thích thành công", nil, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (h *wishlistHandler) GetWishList(ctx *gin.Context) {
	// Get the user id from the context
	user_id := ctx.MustGet("id").(int)
	// Get the order_by query parameter
	order_by := ctx.DefaultQuery("order_by", "default") // Replace "default" with your
	// Perform get wishlist operation
	products, err := h.service.GetWishList(uint(user_id), order_by)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể lấy danh sách yêu thích ", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Lấy danh sách yêu thích thành công", products, nil)
	ctx.JSON(http.StatusOK, successRes)
}
