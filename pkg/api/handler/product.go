package handler

import (
	services "ahava/pkg/service"
	"ahava/pkg/utils/models"
	"ahava/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductHandler interface {
	AddProduct(ctx *gin.Context)
	DeleteProduct(ctx *gin.Context)
	GetProductDetails(ctx *gin.Context)
	ListCategoryProducts(ctx *gin.Context)
	ListFeaturedProducts(ctx *gin.Context)
	// ListProductsForUser(ctx *gin.Context)
	SearchProducts(ctx *gin.Context)
	// GetSearchHistory(ctx *gin.Context)
	UpdateProductImage(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)
	ListProductsForAdmin(ctx *gin.Context)
}

type productHandler struct {
	ProductService services.ProductService
}

func NewProductHandler(service services.ProductService) ProductHandler {
	return &productHandler{
		ProductService: service,
	}
}

func (h *productHandler) AddProduct(ctx *gin.Context) {

	var product models.Product
	category_id, err := strconv.Atoi(ctx.Request.FormValue("category_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "form file error", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	price, err := strconv.Atoi(ctx.Request.FormValue("price"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "form file error", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	stock, err := strconv.Atoi(ctx.Request.FormValue("stock"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "form file error", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	product.CategoryID = uint(category_id)
	product.Name = ctx.Request.FormValue("name")
	product.Size = ctx.Request.FormValue("size")
	product.Price = uint64(price)
	product.Stock = uint(stock)
	product.Description = ctx.Request.FormValue("description")
	product.ShortDescription = ctx.Request.FormValue("short_description")
	product.HowToUse = ctx.Request.FormValue("how_to_use")

	default_image, err := ctx.FormFile("default_image")
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "retrieving image from form error", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "retrieving image from form error", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	images := form.File["images"]

	result, err := h.ProductService.AddProduct(product, default_image, images)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the Product", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added Product", result, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (h *productHandler) DeleteProduct(ctx *gin.Context) {

	product_id, err := strconv.Atoi(ctx.Param("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = h.ProductService.DeleteProduct(uint(product_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully deleted the product", nil, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (h *productHandler) GetProductDetails(ctx *gin.Context) {

	product_id, err := strconv.Atoi(ctx.Query("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	product, err := h.ProductService.GetProductDetails(uint(product_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve product", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Product details retrieved successfully", product, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (h *productHandler) ListCategoryProducts(ctx *gin.Context) {

	category_id, err := strconv.Atoi(ctx.Query("category_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	products, err := h.ProductService.ListCategoryProducts(uint(category_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all home products", products, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (h *productHandler) ListFeaturedProducts(ctx *gin.Context) {

	products, err := h.ProductService.ListFeaturedProducts()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", products, nil)
	ctx.JSON(http.StatusOK, successRes)
}

// func (h *productHandler) ListProductsForUser(ctx *gin.Context) {

// 	id := ctx.MustGet("id")
// 	userID, ok := id.(int)
// 	if !ok {
// 		errorRes := response.ClientResponse(http.StatusForbidden, "problem in identifying user from the context", nil, nil)
// 		ctx.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	products, err := h.ProductService.ListProductsForUser(userID)
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
// 		ctx.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}
// 	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", products, nil)
// 	ctx.JSON(http.StatusOK, successRes)
// }

func (h *productHandler) SearchProducts(ctx *gin.Context) {

	// user_id := ctx.MustGet("id").(int)
	var searchkey models.Search
	if err := ctx.BindJSON(&searchkey); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	results, err := h.ProductService.SearchProducts(searchkey.Key)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve the records", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", results, nil)
	ctx.JSON(http.StatusOK, successRes)
}

// func (h *productHandler) GetSearchHistory(ctx *gin.Context) {

// 	user_id := ctx.MustGet("id").(int)
// 	results, err := h.ProductService.GetSearchHistory(uint(user_id))
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve the records", nil, err.Error())
// 		ctx.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", results, nil)
// 	ctx.JSON(http.StatusOK, successRes)
// }

func (h *productHandler) UpdateProductImage(ctx *gin.Context) {

	product_id, err := strconv.Atoi(ctx.Param("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "parameter problem", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	file, err := ctx.FormFile("default_image")
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "retrieving image from form error", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	results, err := h.ProductService.UpdateProductImage(uint(product_id), file)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not change the image", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully changed image", results, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (h *productHandler) UpdateProduct(ctx *gin.Context) {

	product_id, err := strconv.Atoi(ctx.Param("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "parameter problem", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	var model models.Product
	err = ctx.BindJSON(&model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "parameter problem", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err = validator.New().Struct(model)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errRes)
		return
	}

	results, err := h.ProductService.UpdateProduct(uint(product_id), model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not update product", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated product", results, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (h *productHandler) ListProductsForAdmin(ctx *gin.Context) {
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		offset = 0
	}
	products, err := h.ProductService.ListProductsForAdmin(limit, offset)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", products, nil)
	ctx.JSON(http.StatusOK, successRes)
}
