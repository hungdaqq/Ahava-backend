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
	AddProductImages(ctx *gin.Context)
	DeleteProduct(ctx *gin.Context)
	GetProductDetails(ctx *gin.Context)
	ListCategoryProducts(ctx *gin.Context)
	ListFeaturedProducts(ctx *gin.Context)
	ListAllProducts(ctx *gin.Context)
	SearchProducts(ctx *gin.Context)
	UpdateProductImage(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)
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
	// Bind the request body to the model
	var product models.Product
	err := ctx.BindJSON(&product)
	if err != nil {
		errorRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Validate the model
	result, err := h.ProductService.AddProduct(product)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể thêm sản phẩm", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusCreated, "Thêm sản phẩm thành công", result, nil)
	ctx.JSON(http.StatusCreated, successRes)
}

func (h *productHandler) AddProductImages(ctx *gin.Context) {
	// Get the product id from the context
	product_id, err := strconv.Atoi(ctx.Param("product_id"))
	if err != nil {
		errorRes := response.ClientErrorResponse("Request parameter problem", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Get the default image from the form
	default_image, err := ctx.FormFile("default_image")
	if err != nil {
		errorRes := response.ClientErrorResponse("Tải hình ảnh không thành công", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Get the images from the form
	form, err := ctx.MultipartForm()
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể tải hình ảnh", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	images := form.File["images"]
	// Perform add product images operation
	result, err := h.ProductService.AddProductImages(uint(product_id), default_image, images)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể thêm sản phẩm", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusCreated, "Thêm sản phẩm thành công", result, nil)
	ctx.JSON(http.StatusCreated, successRes)
}

func (h *productHandler) DeleteProduct(ctx *gin.Context) {
	// Get the product id from the context
	product_id, err := strconv.Atoi(ctx.Param("product_id"))
	if err != nil {
		errorRes := response.ClientErrorResponse("Request parameter problem", nil, err)
		ctx.JSON(errorRes.StatusCode, errorRes)
		return
	}
	// Perform delete product operation
	err = h.ProductService.DeleteProduct(uint(product_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Không thể xoá sản phẩm", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Xoá sản phẩm thành công", nil, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (h *productHandler) GetProductDetails(ctx *gin.Context) {
	// Get the product id from the context
	product_id, err := strconv.Atoi(ctx.Query("product_id"))
	if err != nil {
		errorRes := response.ClientErrorResponse("Request query problem", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Perform get product details operation
	product, err := h.ProductService.GetProductDetails(uint(product_id))
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể lấy thông tin sản phẩm", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Lấy thông tin sản phẩm thành công", product, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (h *productHandler) ListCategoryProducts(ctx *gin.Context) {
	// Get the category from the query
	category := ctx.Query("category")
	// Perform list category products operation
	products, err := h.ProductService.ListCategoryProducts(category)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể lấy danh sách sản phẩm", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Lấy danh sách sản phẩm thành công", products, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (h *productHandler) ListFeaturedProducts(ctx *gin.Context) {
	// Perform list featured products operation
	products, err := h.ProductService.ListFeaturedProducts()
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể lấy danh sách sản phẩm", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Lấy danh sách sản phẩm thành công", products, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (h *productHandler) SearchProducts(ctx *gin.Context) {
	// Bind the request body to the model
	var searchkey models.Search
	if err := ctx.BindJSON(&searchkey); err != nil {
		errorRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Perform search products operation
	results, err := h.ProductService.SearchProducts(searchkey.Key)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể lấy danh sách sản phẩm", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Lấy danh sách sản phẩm thành công", results, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (h *productHandler) UpdateProductImage(ctx *gin.Context) {
	// Get the product id from the context
	product_id, err := strconv.Atoi(ctx.Param("product_id"))
	if err != nil {
		errorRes := response.ClientErrorResponse("Request parameter problem", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Get the image from the form
	file, err := ctx.FormFile("default_image")
	if err != nil {
		errorRes := response.ClientErrorResponse("Tải hình ảnh không thành công", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Perform update product image operation
	results, err := h.ProductService.UpdateProductImage(uint(product_id), file)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể cập nhật ảnh sản phẩm", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Cập nhật ảnh sản phẩm thành công", results, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (h *productHandler) UpdateProduct(ctx *gin.Context) {
	// Get the product id from the context
	product_id, err := strconv.Atoi(ctx.Param("product_id"))
	if err != nil {
		errorRes := response.ClientErrorResponse("Request parameter problem", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Bind the request body to the model
	var model models.Product
	err = ctx.BindJSON(&model)
	if err != nil {
		errorRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Validate the model
	err = validator.New().Struct(model)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Constraints not satisfied", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errRes)
		return
	}
	// Perform update product operation
	result, err := h.ProductService.UpdateProduct(uint(product_id), model)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể cập nhật sản phẩm", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Lấy danh sách sản phẩm thành công", result, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (h *productHandler) ListAllProducts(ctx *gin.Context) {
	// Get the limit and offset from the query
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		errorRes := response.ClientErrorResponse("Request query problem", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Get the offset from the query
	offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if err != nil {
		errorRes := response.ClientErrorResponse("Request query problem", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Perform list all products operation
	products, err := h.ProductService.ListAllProducts(limit, offset)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể lấy danh sách sản phẩm", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusOK, "Lấy danh sách sản phẩm thành công", products, nil)
	ctx.JSON(http.StatusOK, successRes)
}
