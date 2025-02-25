package handler

import (
	services "ahava/pkg/service"
	models "ahava/pkg/utils/models"
	response "ahava/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NewsHandler interface {
	AddNews(ctx *gin.Context)
	UpdateNews(ctx *gin.Context)
	DeleteNews(ctx *gin.Context)
	ListAllNews(ctx *gin.Context)
	GetFeaturedNews(ctx *gin.Context)
	GetNewsByID(ctx *gin.Context)
}

type newsHandler struct {
	service services.NewsService
}

func NewNewsHandler(service services.NewsService) NewsHandler {
	return &newsHandler{
		service: service,
	}
}

func (n *newsHandler) AddNews(ctx *gin.Context) {

	var model models.News
	if err := ctx.BindJSON(&model); err != nil {
		errorRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := n.service.AddNews(model)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể thêm tin tức", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusCreated, "Tin tức đã được thêm", result, nil)
	ctx.JSON(http.StatusCreated, successRes)
}

func (n *newsHandler) UpdateNews(ctx *gin.Context) {

	news_id, err := strconv.Atoi(ctx.Param("news_id"))
	if err != nil {
		errorRes := response.ClientErrorResponse("Request query problem", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	var model models.News
	if err := ctx.BindJSON(&model); err != nil {
		errorRes := response.ClientErrorResponse("Fields provided are in wrong format", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := n.service.UpdateNews(uint(news_id), model)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể cập nhật tin tức", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Tin tức đã được cập nhật", result, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (n *newsHandler) DeleteNews(ctx *gin.Context) {

	news_id, err := strconv.Atoi(ctx.Param("news_id"))
	if err != nil {
		errorRes := response.ClientErrorResponse("Request query problem", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = n.service.DeleteNews(uint(news_id))
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể xóa tin tức", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Tin tức đã được xóa", nil, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (n *newsHandler) ListAllNews(ctx *gin.Context) {

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		offset = 0
	}

	result, err := n.service.ListAllNews(limit, offset)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể lấy danh sách tin tức", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Danh sách tin tức", result, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (n *newsHandler) GetFeaturedNews(ctx *gin.Context) {

	result, err := n.service.GetFeaturedNews()
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể lấy danh sách tin tức nổi bật", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Danh sách tin tức nổi bật", result, nil)
	ctx.JSON(http.StatusOK, successRes)
}

func (n *newsHandler) GetNewsByID(ctx *gin.Context) {

	news_id, err := strconv.Atoi(ctx.Param("news_id"))
	if err != nil {
		errorRes := response.ClientErrorResponse("Request query problem", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}

	result, err := n.service.GetNewsByID(uint(news_id))
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể lấy tin tức", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Tin tức", result, nil)
	ctx.JSON(http.StatusOK, successRes)
}
