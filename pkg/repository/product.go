package repository

import (
	"ahava/pkg/domain"
	"ahava/pkg/utils/errors"
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	AddProduct(product models.Products) (models.Products, error)

	UpdateProduct(uint, models.Products) (models.Products, error)
	UpdateProductImage(uint, string) (models.Products, error)

	DeleteProduct(product_id uint) error

	GetProductDetails(product_id uint) (models.Products, error)
	ListProducts(limit, offset int) (models.ListProducts, error)
	ListCategoryProducts(category_id uint) ([]models.Products, error)
	ListFeaturedProducts() ([]models.Products, error)

	SearchProducts(key string) ([]models.Products, error)
}

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) *productRepository {
	return &productRepository{
		DB: DB,
	}
}

func (i *productRepository) AddProduct(product models.Products) (models.Products, error) {

	var addProduct models.Products

	err := i.DB.Create(&domain.Products{
		Name:             product.Name,
		CategoryID:       product.CategoryID,
		Price:            product.Price,
		Size:             product.Size,
		Stock:            product.Stock,
		DefaultImage:     product.DefaultImage,
		Images:           product.Images,
		ShortDescription: product.ShortDescription,
		Description:      product.Description,
		HowToUse:         product.HowToUse,
	}).Scan(&addProduct).Error
	if err != nil {
		return models.Products{}, err
	}

	return addProduct, nil
}

func (i *productRepository) DeleteProduct(product_id uint) error {

	result := i.DB.Exec("DELETE FROM products WHERE id = ?", product_id)

	if result.RowsAffected < 1 {
		return errors.ErrEntityNotFound
	}

	return nil
}

func (i *productRepository) GetProductDetails(product_id uint) (models.Products, error) {

	var product models.Products

	err := i.DB.Raw(`SELECT * FROM products WHERE products.id = ?`,
		product_id).Scan(&product).Error
	if err != nil {
		return models.Products{}, errors.ErrEntityNotFound
	}

	return product, nil
}

func (i *productRepository) ListProducts(limit, offset int) (models.ListProducts, error) {

	var listProducts models.ListProducts
	var productDetails []models.Products
	var total int64

	query := i.DB.Model(&models.Products{})
	// Get the total count of records
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

func (i *productRepository) ListCategoryProducts(category_id uint) ([]models.Products, error) {

	var products []models.Products
	err := i.DB.Raw("SELECT * FROM products WHERE category_id=$1", category_id).
		Scan(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (i *productRepository) ListFeaturedProducts() ([]models.Products, error) {

	var products []models.Products
	err := i.DB.Raw("SELECT * FROM products WHERE is_featured=true").Scan(&products).Error
	if err != nil {
		return []models.Products{}, err
	}

	return products, nil
}

// func (i *productRepository) CheckStock(product_id uint) (int, error) {
// 	var stock int
// 	if err := i.DB.Raw("SELECT stock FROM products WHERE id=$1", product_id).Scan(&stock).Error; err != nil {
// 		return 0, err
// 	}
// 	return stock, nil
// }

func (i *productRepository) SearchProducts(key string) ([]models.Products, error) {

	var productDetails []models.Products

	query := `SELECT i.* FROM products i LEFT JOIN categories c ON i.category_id = c.id 
				WHERE i.name ILIKE '%' || ? || '%' OR c.name ILIKE '%' || ? || '%'`
	if err := i.DB.Raw(query, key, key).Scan(&productDetails).Error; err != nil {
		return []models.Products{}, err
	}

	return productDetails, nil
}

// func (i *productRepository) SaveSearchHistory(user_id uint, key string) error {

// 	err := i.DB.Exec(`INSERT INTO search_histories (user_id,search_key) VALUES (?,?)`,
// 		user_id, key).Error
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (i *productRepository) GetSearchHistory(user_id uint) ([]models.SearchHistory, error) {

// 	var searchHistory []models.SearchHistory

// 	err := i.DB.Raw(`SELECT * FROM search_histories WHERE user_id=$1 ORDER BY created_at DESC`,
// 		user_id).Scan(&searchHistory).Error
// 	if err != nil {
// 		return []models.SearchHistory{}, err
// 	}

// 	return searchHistory, nil
// }

func (i *productRepository) UpdateProductImage(product_id uint, url string) (models.Products, error) {

	var updateProduct models.Products

	// Use GORM's Update method to update the product image
	result := i.DB.Model(&domain.Products{}).Where("id = ?", product_id).Update("image", url).Scan(&updateProduct)

	// Check for errors in the update process
	if result.Error != nil {
		return models.Products{}, result.Error
	}

	// Check if any row was affected (to handle case where no rows matched)
	if result.RowsAffected == 0 {
		return models.Products{}, errors.ErrEntityNotFound
	}

	return updateProduct, nil
}

func (i *productRepository) UpdateProduct(product_id uint, model models.Products) (models.Products, error) {

	var updateProduct models.Products

	result := i.DB.Model(&domain.Products{}).Where("id = ?", product_id).Updates(
		domain.Products{
			Name:             model.Name,
			CategoryID:       model.CategoryID,
			Price:            model.Price,
			Size:             model.Size,
			Stock:            model.Stock,
			Description:      model.Description,
			ShortDescription: model.ShortDescription,
			IsFeatured:       model.IsFeatured,
			HowToUse:         model.HowToUse,
		}).Scan(&updateProduct)

	if result.Error != nil {
		return models.Products{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.Products{}, errors.ErrEntityNotFound
	}

	return updateProduct, nil
}
