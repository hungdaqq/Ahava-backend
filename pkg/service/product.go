package service

import (
	helper "ahava/pkg/helper"
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
)

type ProductService interface {
	AddProduct(product models.Product) (models.Product, error)
	UpdateProduct(uint, models.Product) (models.Product, error)
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

func (i *productService) AddProduct(p models.Product) (models.Product, error) {
	// Add product
	product, err := i.repository.AddProduct(p)
	if err != nil {
		return models.Product{}, err
	}
	// Add product price
	prices := []models.Price{}
	for _, pr := range p.Price {
		price, err := i.repository.AddProductPrice(product.ID, pr)
		if err != nil {
			return models.Product{}, err
		}
		prices = append(prices, price)
	}
	// Assign the price to the product
	product.Price = prices
	// Return the product
	return product, nil
}

func (i *productService) UpdateProduct(product_id uint, p models.Product) (models.Product, error) {
	// Update product details
	product, err := i.repository.UpdateProduct(product_id, p)
	if err != nil {
		return models.Product{}, err
	}
	// Get old prices
	old_prices, err := i.repository.GetProductPrice(product_id)
	if err != nil {
		return models.Product{}, err
	}
	// Initialize slices to track prices
	var prices []models.Price
	updatedPriceIDs := make(map[uint]struct{}) // Use a map to track updated prices
	// Flag to check if all prices are new (ID == 0)
	allNewPrices := true
	// Handle adding new prices and updating existing ones
	for _, pr := range p.Price {
		// If the price is new (ID == 0), add it
		if pr.ID == 0 {
			price, err := i.repository.AddProductPrice(product_id, pr)
			if err != nil {
				return models.Product{}, err
			}
			prices = append(prices, price)
			updatedPriceIDs[price.ID] = struct{}{}
		} else {
			// If the price already exists, update it
			price, err := i.repository.UpdateProductPrice(product_id, pr.ID, pr)
			if err != nil {
				return models.Product{}, err
			}
			prices = append(prices, price)
			updatedPriceIDs[pr.ID] = struct{}{}
			allNewPrices = false // There's at least one existing price, so not all are new
		}
	}
	// If not all prices are new, delete obsolete prices
	if !allNewPrices {
		// Identify which old prices are no longer in the updated list (obsolete prices)
		for _, oldPrice := range old_prices {
			// If this price ID is not in the updated list, delete it
			if _, exists := updatedPriceIDs[oldPrice.ID]; !exists {
				err := i.repository.DeleteProductPrice(product_id, oldPrice.ID)
				if err != nil {
					return models.Product{}, err
				}
			}
		}
	}
	// Assign the updated prices to the product
	product.Price = prices
	// Return the updated product
	return product, nil
}

func (i *productService) DeleteProduct(product_id uint) error {
	// Delete product
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
