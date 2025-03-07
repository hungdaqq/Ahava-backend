package handler

import (
	services "ahava/pkg/service"
	"ahava/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadHandler interface {
	FileUpload(ctx *gin.Context)
}

type uploadHandler struct {
	uploadService services.UploadService
}

func NewUploadHandler(service services.UploadService) UploadHandler {
	return &uploadHandler{
		uploadService: service,
	}
}

func (h *uploadHandler) FileUpload(ctx *gin.Context) {
	// Get the files from the form
	form, err := ctx.MultipartForm()
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể tải hình ảnh", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	files := form.File["files"]
	// Perform add product images operation
	result, err := h.uploadService.FileUpload(files)
	if err != nil {
		errorRes := response.ClientErrorResponse("Không thể thêm hình ảnh", nil, err)
		ctx.JSON(http.StatusBadRequest, errorRes)
		return
	}
	// Return the response
	successRes := response.ClientResponse(http.StatusCreated, "Thêm hình ảnh thành công", result, nil)
	ctx.JSON(http.StatusCreated, successRes)
}
