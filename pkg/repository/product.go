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

	AddProductPrice(product_id uint, price []models.Price) ([]models.Price, error)
	UpdateProductPrice(product_id uint, price []models.Price) ([]models.Price, error)
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
		IsFeatured:       p.IsFeatured,
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
		IsFeatured:       product.IsFeatured,
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

	var product domain.Product

	err := r.DB.First(&product, product_id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Product{}, models.ErrEntityNotFound
		}
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
		IsFeatured:       product.IsFeatured,
	}, nil
}

func (r *productRepository) ListAllProducts(limit, offset int) (models.ListProducts, error) {

	var listProducts models.ListProducts
	var productDetails []models.Product
	var total int64

	query := r.DB.Model(&domain.Product{})
	if err := query.Count(&total).Error; err != nil {
		return models.ListProducts{}, err
	}
	if err := query.Offset(offset).Limit(limit).Find(&productDetails).Error; err != nil {
		return models.ListProducts{}, err
	}

	listProducts.Products = productDetails
	listProducts.Total = total
	listProducts.Limit = limit
	listProducts.Offset = offset

	return listProducts, nil
}

func (r *productRepository) ListCategoryProducts(category string) ([]models.Product, error) {

	var products []models.Product

	err := r.DB.Model(&domain.Product{}).
		Where("category = ?", category).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) ListFeaturedProducts() ([]models.Product, error) {

	var products []models.Product

	err := r.DB.Model(&domain.Product{}).
		Where("is_featured = true").
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) SearchProducts(key string) ([]models.Product, error) {

	var products []models.Product

	err := r.DB.Model(&domain.Product{}).
		Where("name ILIKE ? OR category ILIKE ?", "%"+key+"%", "%"+key+"%").
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) UpdateProduct(product_id uint, p models.Product) (models.Product, error) {

	var product models.Product

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
				IsFeatured:       p.IsFeatured,
				HowToUse:         p.HowToUse,
			}).
		Scan(&product)
	if result.Error != nil {
		return models.Product{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.Product{}, models.ErrEntityNotFound
	}

	return product, nil
}

func (r *productRepository) UpdateProductPrice(product_id uint, price []models.Price) ([]models.Price, error) {

	var updatePrice []models.Price
	for _, p := range price {
		result := r.DB.Model(&domain.Price{}).
			Where("product_id = ? AND size = ?", product_id, p.Size).
			Updates(
				domain.Price{
					OriginalPrice: p.OriginalPrice,
					DiscountPrice: p.DiscountPrice,
				}).
			Scan(&p)
		if result.Error != nil {
			return nil, result.Error
		}
		if result.RowsAffected == 0 {
			return nil, models.ErrEntityNotFound
		}

		updatePrice = append(updatePrice, models.Price{
			Size:          p.Size,
			OriginalPrice: p.OriginalPrice,
			DiscountPrice: p.DiscountPrice,
		})
	}

	return updatePrice, nil
}

func (r *productRepository) AddProductPrice(product_id uint, price []models.Price) ([]models.Price, error) {

	var updatePrice []models.Price

	for _, p := range price {
		price := domain.Price{
			ProductID:     product_id,
			Size:          p.Size,
			OriginalPrice: p.OriginalPrice,
			DiscountPrice: p.DiscountPrice,
		}
		if err := r.DB.Create(&price).Error; err != nil {
			return nil, err
		}
		updatePrice = append(updatePrice, models.Price{
			Size:          price.Size,
			OriginalPrice: price.OriginalPrice,
			DiscountPrice: price.DiscountPrice,
		})
	}

	return updatePrice, nil
}
func (r *productRepository) GetProductPrice(product_id uint) ([]models.Price, error) {

	var prices []models.Price

	err := r.DB.Model(&domain.Price{}).
		Where("product_id = ?", product_id).
		Scan(&prices).Error
	if err != nil {
		return nil, err
	}

	return prices, nil
}
