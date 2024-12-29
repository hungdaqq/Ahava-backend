package handler

import (
	services "ahava/pkg/usecase"
	"ahava/pkg/utils/models"
	"ahava/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductHandler struct {
	ProductUseCase services.ProductUseCase
}

func NewProductHandler(usecase services.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		ProductUseCase: usecase,
	}
}

func (i *ProductHandler) AddProduct(c *gin.Context) {

	var product models.Products
	category_id, err := strconv.Atoi(c.Request.FormValue("category_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "form file error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	price, err := strconv.Atoi(c.Request.FormValue("price"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "form file error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	stock, err := strconv.Atoi(c.Request.FormValue("stock"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "form file error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	product.CategoryID = uint(category_id)
	product.Name = c.Request.FormValue("name")
	product.Size = c.Request.FormValue("size")
	product.Price = uint64(price)
	product.Stock = uint(stock)
	product.Description = c.Request.FormValue("description")
	product.ShortDescription = c.Request.FormValue("short_description")
	product.HowToUse = c.Request.FormValue("how_to_use")

	default_image, err := c.FormFile("default_image")
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "retrieving image from form error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "retrieving image from form error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	images := form.File["images"]

	result, err := i.ProductUseCase.AddProduct(product, default_image, images)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the Product", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added Product", result, nil)
	c.JSON(http.StatusOK, successRes)

}

func (i *ProductHandler) DeleteProduct(c *gin.Context) {

	product_id, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = i.ProductUseCase.DeleteProduct(uint(product_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully deleted the product", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

func (i *ProductHandler) GetProductDetails(c *gin.Context) {

	product_id, err := strconv.Atoi(c.Query("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	product, err := i.ProductUseCase.GetProductDetails(uint(product_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve product", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Product details retrieved successfully", product, nil)
	c.JSON(http.StatusOK, successRes)
}

func (i *ProductHandler) ListCategoryProducts(c *gin.Context) {

	category_id, err := strconv.Atoi(c.Query("category_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	products, err := i.ProductUseCase.ListCategoryProducts(uint(category_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all home products", products, nil)
	c.JSON(http.StatusOK, successRes)
}

func (i *ProductHandler) ListFeaturedProducts(c *gin.Context) {

	products, err := i.ProductUseCase.ListFeaturedProducts()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", products, nil)
	c.JSON(http.StatusOK, successRes)
}

// func (i *ProductHandler) ListProductsForUser(c *gin.Context) {

// 	id := c.MustGet("id")
// 	userID, ok := id.(int)
// 	if !ok {
// 		errorRes := response.ClientResponse(http.StatusForbidden, "problem in identifying user from the context", nil, nil)
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	products, err := i.ProductUseCase.ListProductsForUser(userID)
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}
// 	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", products, nil)
// 	c.JSON(http.StatusOK, successRes)
// }

func (i *ProductHandler) SearchProducts(c *gin.Context) {

	// user_id := c.MustGet("id").(int)
	var searchkey models.Search
	if err := c.BindJSON(&searchkey); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	results, err := i.ProductUseCase.SearchProducts(searchkey.Key)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve the records", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", results, nil)
	c.JSON(http.StatusOK, successRes)
}

// func (i *ProductHandler) GetSearchHistory(c *gin.Context) {

// 	user_id := c.MustGet("id").(int)
// 	results, err := i.ProductUseCase.GetSearchHistory(uint(user_id))
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve the records", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", results, nil)
// 	c.JSON(http.StatusOK, successRes)
// }

func (i *ProductHandler) UpdateProductImage(c *gin.Context) {

	product_id, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "parameter problem", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	file, err := c.FormFile("default_image")
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "retrieving image from form error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	results, err := i.ProductUseCase.UpdateProductImage(uint(product_id), file)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not change the image", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully changed image", results, nil)
	c.JSON(http.StatusOK, successRes)

}

func (i *ProductHandler) UpdateProduct(c *gin.Context) {

	product_id, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "parameter problem", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	var model models.Products
	err = c.BindJSON(&model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "parameter problem", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err = validator.New().Struct(model)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	results, err := i.ProductUseCase.UpdateProduct(uint(product_id), model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not update product", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated product", results, nil)
	c.JSON(http.StatusOK, successRes)

}

func (i *ProductHandler) ListProductsForAdmin(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}
	products, err := i.ProductUseCase.ListProductsForAdmin(limit, offset)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "could not retrieve records", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Successfully got all records", products, nil)
	c.JSON(http.StatusOK, successRes)
}
