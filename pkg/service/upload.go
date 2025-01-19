package service

import (
	helper "ahava/pkg/helper"
	"mime/multipart"
)

type UploadService interface {
	FileUpload(file *multipart.FileHeader) (string, error)
}

type uploadService struct {
	helper helper.Helper
}

func NewUploadService(
	h helper.Helper) UploadService {

	return &uploadService{helper: h}
}

func (u *uploadService) FileUpload(file *multipart.FileHeader) (string, error) {
	// Upload file to S3
	file_name, err := u.helper.AddFileToS3(file, "ahava")
	if err != nil {
		return "", err
	}

	return file_name, nil
}
