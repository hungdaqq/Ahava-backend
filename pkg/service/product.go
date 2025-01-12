package service

import (
	helper "ahava/pkg/helper"
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
	"mime/multipart"

	"github.com/lib/pq"
)

type ProductService interface {
	AddProduct(product models.Product) (models.Product, error)
	AddProductImages(product_id uint, default_image *multipart.FileHeader, images []*multipart.FileHeader) (models.Product, error)

	UpdateProduct(uint, models.Product) (models.Product, error)
	UpdateProductImage(product_id uint, file *multipart.FileHeader) (models.Product, error)

	DeleteProduct(product_id uint) error

	GetProductDetails(product_id uint) (models.Product, error)
	ListAllProducts(limit, offest int) (models.ListProducts, error)
	ListCategoryProducts(category string) ([]models.Product, error)
	ListFeaturedProducts() ([]models.Product, error)
	SearchProducts(key string) ([]models.Product, error)
}

type productService struct {
	repository repository.ProductRepository
	helper     helper.Helper
}

func NewProductService(
	repo repository.ProductRepository,
	h helper.Helper,
) ProductService {
	return &productService{
		repository: repo,
		helper:     h,
	}
}

func (i *productService) AddProduct(product models.Product) (models.Product, error) {

	addProduct, err := i.repository.AddProduct(product)
	if err != nil {
		return models.Product{}, err
	}

	addPrice, err := i.repository.AddProductPrice(addProduct.ID, product.Price)
	if err != nil {
		return models.Product{}, err
	}

	addProduct.Price = addPrice

	return addProduct, nil
}

func (i *productService) AddProductImages(product_id uint, default_image *multipart.FileHeader, images []*multipart.FileHeader) (models.Product, error) {

	default_image_url, err := i.helper.AddImageToS3(default_image, "ahava")
	if err != nil {
		return models.Product{}, err
	}

	var urls []string

	for _, image := range images {
		url, err := i.helper.AddImageToS3(image, "ahava")
		if err != nil {
			return models.Product{}, err
		}

		urls = append(urls, url)
	}

	result, err := i.repository.AddProductImages(product_id, default_image_url, pq.StringArray(urls))
	if err != nil {
		return models.Product{}, err
	}
	return result, nil
}

func (i *productService) UpdateProductImage(id uint, file *multipart.FileHeader) (models.Product, error) {

	url, err := i.helper.AddImageToS3(file, "ahava")
	if err != nil {
		return models.Product{}, err
	}

	//send the url and save it in database
	result, err := i.repository.UpdateProductImage(id, url)
	if err != nil {
		return models.Product{}, err
	}

	return result, nil

}

func (i *productService) UpdateProduct(product_id uint, model models.Product) (models.Product, error) {

	updateProduct, err := i.repository.UpdateProduct(product_id, model)
	if err != nil {
		return models.Product{}, err
	}

	updatePrice, err := i.repository.UpdateProductPrice(product_id, model.Price)
	if err != nil {
		return models.Product{}, err
	}

	updateProduct.Price = updatePrice

	return updateProduct, nil
}

func (i *productService) DeleteProduct(product_id uint) error {

	err := i.repository.DeleteProduct(product_id)
	if err != nil {
		return err
	}

	return nil
}

func (i *productService) GetProductDetails(product_id uint) (models.Product, error) {

	product, err := i.repository.GetProductDetails(product_id)
	if err != nil {
		return models.Product{}, err
	}

	price, err := i.repository.GetProductPrice(product_id)
	if err != nil {
		return models.Product{}, err
	}

	product.Price = price

	return product, nil
}

func (i *productService) ListCategoryProducts(category string) ([]models.Product, error) {

	products, err := i.repository.ListCategoryProducts(category)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (i *productService) ListFeaturedProducts() ([]models.Product, error) {

	products, err := i.repository.ListFeaturedProducts()
	if err != nil {
		return []models.Product{}, err
	}

	for idx := range products {
		price, err := i.repository.GetProductPrice(products[idx].ID)
		if err != nil {
			return []models.Product{}, err
		}
		products[idx].Price = price
	}

	return products, nil
}

func (i *productService) ListAllProducts(limit, offset int) (models.ListProducts, error) {

	products, err := i.repository.ListAllProducts(limit, offset)
	if err != nil {
		return models.ListProducts{}, err
	}

	for idx := range products.Products {
		price, err := i.repository.GetProductPrice(products.Products[idx].ID)
		if err != nil {
			return models.ListProducts{}, err
		}
		products.Products[idx].Price = price
	}

	return products, nil
}

func (i *productService) SearchProducts(key string) ([]models.Product, error) {

	products, err := i.repository.SearchProducts(key)
	if err != nil {
		return []models.Product{}, err
	}

	for idx := range products {
		price, err := i.repository.GetProductPrice(products[idx].ID)
		if err != nil {
			return []models.Product{}, err
		}
		products[idx].Price = price
	}

	return products, nil
}
