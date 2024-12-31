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

type CategoryHandler interface {
	AddCategory(ctx *gin.Context)
	UpdateCategory(ctx *gin.Context)
	DeleteCategory(ctx *gin.Context)
	GetCategory(ctx *gin.Context)
	// GetProductDetailsInACategory(ctx *gin.Context)
	// GetBannersForUsers(ctx *gin.Context)
}

type categoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(service services.CategoryService) CategoryHandler {
	return &categoryHandler{
		categoryService: service,
	}
}

func (h *categoryHandler) AddCategory(ctx *gin.Context) {

	var model models.Category
	err := ctx.BindJSON(&model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err = validator.New().Struct(model)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Constraints not satisfied", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errRes)
		return
	}

	categoryResponse, err := h.categoryService.AddCategory(model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the Category", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added the Category", categoryResponse, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (h *categoryHandler) UpdateCategory(ctx *gin.Context) {

	category_id, err := strconv.Atoi(ctx.Param("category_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Parameter problem", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	var p models.UpdateCategory
	if err := ctx.BindJSON(&p); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	a, err := h.categoryService.UpdateCategory(uint(category_id), p.Name, p.Description)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not update the Category", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated the Category", a, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (h *categoryHandler) DeleteCategory(ctx *gin.Context) {

	category_id, err := strconv.Atoi(ctx.Param("category_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "parameter problem", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = h.categoryService.DeleteCategory(uint(category_id))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully deleted the Category", nil, nil)
	ctx.JSON(http.StatusOK, successRes)

}

func (h *categoryHandler) GetCategory(ctx *gin.Context) {

	categories, err := h.categoryService.GetCategories()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all categories", categories, nil)
	ctx.JSON(http.StatusOK, successRes)

}

// func (h *categoryHandler) GetProductDetailsInACategory(ctx *gin.Context) {

// 	id, err := strconv.Atoi(ctx.Query("id"))
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
// 		ctx.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	products, err := h.categoryService.GetProductDetailsInACategory(id)
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusInternalServerError, "error in fetching data", nil, err.Error())
// 		ctx.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	successRes := response.ClientResponse(http.StatusOK, "Successfully got all categories", products, nil)
// 	ctx.JSON(http.StatusOK, successRes)

// }

// func (h *categoryHandler) GetBannersForUsers(ctx *gin.Context) {

// 	banners, err := h.categoryService.GetBannersForUsers()
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusInternalServerError, "error in fetching data", nil, err.Error())
// 		ctx.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	successRes := response.ClientResponse(http.StatusOK, "Successfully got all banners", banners, nil)
// 	ctx.JSON(http.StatusOK, successRes)

// }
