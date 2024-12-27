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

type CategoryHandler struct {
	CategoryUseCase services.CategoryUseCase
}

func NewCategoryHandler(usecase services.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		CategoryUseCase: usecase,
	}
}

func (Cat *CategoryHandler) AddCategory(c *gin.Context) {

	var model models.Category
	err := c.BindJSON(&model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err = validator.New().Struct(model)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "Constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	categoryResponse, err := Cat.CategoryUseCase.AddCategory(model)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not add the Category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully added the Category", categoryResponse, nil)
	c.JSON(http.StatusOK, successRes)

}

func (Cat *CategoryHandler) UpdateCategory(c *gin.Context) {

	categoryID, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Parameter problem", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	var p models.UpdateCategory
	if err := c.BindJSON(&p); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	a, err := Cat.CategoryUseCase.UpdateCategory(uint(categoryID), p.CategoryName, p.Description)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not update the Category", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully updated the Category", a, nil)
	c.JSON(http.StatusOK, successRes)

}

func (Cat *CategoryHandler) DeleteCategory(c *gin.Context) {

	categoryID, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "parameter problem", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = Cat.CategoryUseCase.DeleteCategory(uint(categoryID))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully deleted the Category", nil, nil)
	c.JSON(http.StatusOK, successRes)

}

func (Cat *CategoryHandler) GetCategory(c *gin.Context) {

	categories, err := Cat.CategoryUseCase.GetCategories()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully got all categories", categories, nil)
	c.JSON(http.StatusOK, successRes)

}

// func (Cat *CategoryHandler) GetProductDetailsInACategory(c *gin.Context) {

// 	id, err := strconv.Atoi(c.Query("id"))
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	products, err := Cat.CategoryUseCase.GetProductDetailsInACategory(id)
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusInternalServerError, "error in fetching data", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	successRes := response.ClientResponse(http.StatusOK, "Successfully got all categories", products, nil)
// 	c.JSON(http.StatusOK, successRes)

// }

// func (Cat *CategoryHandler) GetBannersForUsers(c *gin.Context) {

// 	banners, err := Cat.CategoryUseCase.GetBannersForUsers()
// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusInternalServerError, "error in fetching data", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	successRes := response.ClientResponse(http.StatusOK, "Successfully got all banners", banners, nil)
// 	c.JSON(http.StatusOK, successRes)

// }
