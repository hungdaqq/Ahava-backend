package service

import (
	helper "ahava/pkg/helper"
	"mime/multipart"
)

type UploadService interface {
	FileUpload(files []*multipart.FileHeader) ([]string, error)
}

type uploadService struct {
	helper helper.Helper
}

func NewUploadService(
	h helper.Helper) UploadService {

	return &uploadService{helper: h}
}

func (s *uploadService) FileUpload(files []*multipart.FileHeader) ([]string, error) {
	// Upload files to S3
	var urls []string
	for _, file := range files {
		url, err := s.helper.AddFileToS3(file, "ahava")
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}
