package repository

import (
	"ahava/pkg/domain"
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	AddProduct(product models.Product) (models.Product, error)

	UpdateProduct(product_id uint, product models.Product) (models.Product, error)

	DeleteProduct(product_id uint) error

	GetProductDetails(product_id uint) (models.Product, error)
	ListAllProducts(limit, offset int) (models.ListProducts, error)
	ListCategoryProducts(category string) ([]models.Product, error)
	ListFeaturedProducts() ([]models.Product, error)

	SearchProducts(key string) ([]models.Product, error)

	GetProductPrice(product_id uint) ([]models.Price, error)

	AddProductPrice(product_id uint, price models.Price) (models.Price, error)
	UpdateProductPrice(product_id, price_id uint, price models.Price) (models.Price, error)
	DeleteProductPrice(product_id, price_id uint) error
}

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) ProductRepository {
	return &productRepository{DB}
}

func (r *productRepository) AddProduct(p models.Product) (models.Product, error) {

	product := domain.Product{
		Name:             p.Name,
		Code:             p.Code,
		Category:         p.Category,
		DefaultImage:     p.DefaultImage,
		Images:           p.Images,
		Stock:            p.Stock,
		ShortDescription: p.ShortDescription,
		Description:      p.Description,
		HowToUse:         p.HowToUse,
		IsFeatured:       &p.IsFeatured,
	}

	if err := r.DB.Create(&product).Error; err != nil {
		return models.Product{}, err
	}

	return models.Product{
		ID:               product.ID,
		Name:             product.Name,
		Code:             product.Code,
		Category:         product.Category,
		DefaultImage:     product.DefaultImage,
		Images:           product.Images,
		Stock:            product.Stock,
		ShortDescription: product.ShortDescription,
		Description:      product.Description,
		HowToUse:         product.HowToUse,
		IsFeatured:       *product.IsFeatured,
	}, nil
}

func (r *productRepository) DeleteProduct(product_id uint) error {

	result := r.DB.Delete(&domain.Product{}, product_id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrEntityNotFound
	}

	return nil
}

func (r *productRepository) GetProductDetails(product_id uint) (models.Product, error) {
	// Define the product
	var product domain.Product
	// Query to get the product details
	err := r.DB.First(&product, product_id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Product{}, models.ErrEntityNotFound
		}
		return models.Product{}, err
	}
	// Return the product details
	return models.Product{
		ID:               product.ID,
		Name:             product.Name,
		Code:             product.Code,
		Category:         product.Category,
		DefaultImage:     product.DefaultImage,
		Images:           product.Images,
		Stock:            product.Stock,
		ShortDescription: product.ShortDescription,
		Description:      product.Description,
		HowToUse:         product.HowToUse,
		IsFeatured:       *product.IsFeatured,
	}, nil
}

func (r *productRepository) ListAllProducts(limit, offset int) (models.ListProducts, error) {
	// Define list of products and product details
	var listProducts models.ListProducts
	var productDetails []models.Product
	var total int64
	// Define the query
	query := r.DB.Model(&domain.Product{})
	if err := query.Count(&total).Error; err != nil {
		return models.ListProducts{}, err
	}
	if err := query.Offset(offset).Limit(limit).Find(&productDetails).Error; err != nil {
		return models.ListProducts{}, err
	}
	// Assign the values to the listProducts
	listProducts.Products = productDetails
	listProducts.Total = total
	listProducts.Limit = limit
	listProducts.Offset = offset
	// Return the list of products
	return listProducts, nil
}

func (r *productRepository) ListCategoryProducts(category string) ([]models.Product, error) {
	// Define list of products and product details
	var products []models.Product
	// Query to get the products based on the category
	err := r.DB.Model(&domain.Product{}).
		Where("category = ?", category).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	// Return the list of products
	return products, nil
}

func (r *productRepository) ListFeaturedProducts() ([]models.Product, error) {
	// Define list of products and product details
	var products []models.Product
	// Query to get the featured products
	err := r.DB.Model(&domain.Product{}).
		Where("is_featured = true").
		Scan(&products).Error
	if err != nil {
		return nil, err
	}
	// Return the list of products
	return products, nil
}

func (r *productRepository) SearchProducts(key string) ([]models.Product, error) {
	// Define list of products and product details
	var products []models.Product
	// Query to search the products based on the key
	err := r.DB.Model(&domain.Product{}).
		Where("name ILIKE ? OR category ILIKE ?", "%"+key+"%", "%"+key+"%").
		Scan(&products).Error
	if err != nil {
		return nil, err
	}
	// Return the list of products
	return products, nil
}

func (r *productRepository) UpdateProduct(product_id uint, p models.Product) (models.Product, error) {
	// Define the product
	var product models.Product
	// Update the product details
	result := r.DB.Model(&domain.Product{}).
		Where("id = ?", product_id).
		Updates(
			domain.Product{
				Name:             p.Name,
				Code:             p.Code,
				Category:         p.Category,
				DefaultImage:     p.DefaultImage,
				Images:           p.Images,
				Stock:            p.Stock,
				Description:      p.Description,
				ShortDescription: p.ShortDescription,
				IsFeatured:       &p.IsFeatured,
				HowToUse:         p.HowToUse,
			}).
		Scan(&product)
	if result.Error != nil {
		return models.Product{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.Product{}, models.ErrEntityNotFound
	}
	// Return the updated product details
	return product, nil
}

func (r *productRepository) UpdateProductPrice(product_id, price_id uint, p models.Price) (models.Price, error) {
	// Define the price
	var price models.Price
	// Update the price
	result := r.DB.Model(&domain.Price{}).
		Where("product_id=? AND id=?", product_id, price_id).
		Updates(
			domain.Price{
				Size:          p.Size,
				Image:         p.Image,
				OriginalPrice: p.OriginalPrice,
				DiscountPrice: p.DiscountPrice,
			}).
		Scan(&price)
	if result.Error != nil {
		return models.Price{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.Price{}, models.ErrEntityNotFound
	}
	// Return the updated price
	return price, nil
}

func (r *productRepository) AddProductPrice(product_id uint, p models.Price) (models.Price, error) {
	// Define the price
	price := domain.Price{
		ProductID:     product_id,
		Size:          p.Size,
		Image:         p.Image,
		OriginalPrice: p.OriginalPrice,
		DiscountPrice: p.DiscountPrice,
	}
	// Create the price
	if err := r.DB.Create(&price).Error; err != nil {
		return models.Price{}, err
	}
	// Return the price
	return models.Price{
		ID:            price.ID,
		Size:          price.Size,
		Image:         price.Image,
		OriginalPrice: price.OriginalPrice,
		DiscountPrice: price.DiscountPrice,
	}, nil
}

func (r *productRepository) GetProductPrice(product_id uint) ([]models.Price, error) {
	// Define the price
	var prices []models.Price
	// Query to get the price details
	err := r.DB.Model(&domain.Price{}).
		Where("product_id = ?", product_id).
		Scan(&prices).Error
	if err != nil {
		return nil, err
	}
	// Return the price details
	return prices, nil
}

func (r *productRepository) DeleteProductPrice(product_id, price_id uint) error {
	// Query to delete the price
	result := r.DB.Delete(&domain.Price{}, price_id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrEntityNotFound
	}
	// Return the price
	return nil
}
