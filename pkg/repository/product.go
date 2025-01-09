package repository

import (
	"ahava/pkg/domain"
	"ahava/pkg/utils/errors"
	"ahava/pkg/utils/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	AddProduct(product models.Product) (models.Product, error)

	UpdateProduct(uint, models.Product) (models.Product, error)
	UpdateProductImage(uint, string) (models.Product, error)

	DeleteProduct(product_id uint) error

	GetProductDetails(product_id uint) (models.Product, error)
	ListProducts(limit, offset int) (models.ListProducts, error)
	ListCategoryProducts(category_id uint) ([]models.Product, error)
	ListFeaturedProducts() ([]models.Product, error)

	SearchProducts(key string) ([]models.Product, error)
}

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) ProductRepository {
	return &productRepository{
		DB: DB,
	}
}

func (r *productRepository) AddProduct(product models.Product) (models.Product, error) {

	var addProduct domain.Product

	err := r.DB.Create(&domain.Product{
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
		return models.Product{}, err
	}

	return models.Product{
		ID:               addProduct.ID,
		Name:             addProduct.Name,
		CategoryID:       addProduct.CategoryID,
		Price:            addProduct.Price,
		Size:             addProduct.Size,
		Stock:            addProduct.Stock,
		DefaultImage:     addProduct.DefaultImage,
		Images:           addProduct.Images,
		ShortDescription: addProduct.ShortDescription,
		Description:      addProduct.Description,
		HowToUse:         addProduct.HowToUse,
	}, nil
}

func (r *productRepository) DeleteProduct(product_id uint) error {

	result := r.DB.Exec("DELETE FROM products WHERE id = ?", product_id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.ErrEntityNotFound
	}

	return nil
}

func (r *productRepository) GetProductDetails(product_id uint) (models.Product, error) {

	var product domain.Product

	err := r.DB.Raw(`SELECT * FROM products WHERE products.id = ?`,
		product_id).Scan(&product).Error
	if err != nil {
		return models.Product{}, errors.ErrEntityNotFound
	}

	return models.Product{
		ID:               product.ID,
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
		IsFeatured:       product.IsFeatured,
	}, nil
}

func (r *productRepository) ListProducts(limit, offset int) (models.ListProducts, error) {

	var listProducts models.ListProducts
	var productDetails []models.Product
	var total int64

	query := r.DB.Model(&models.Product{})
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

func (r *productRepository) ListCategoryProducts(category_id uint) ([]models.Product, error) {

	var products []models.Product
	err := r.DB.Raw("SELECT * FROM products WHERE category_id=$1", category_id).
		Scan(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil

}
func (r *productRepository) ListFeaturedProducts() ([]models.Product, error) {

	var products []models.Product
	err := r.DB.Raw("SELECT * FROM products WHERE is_featured=true").Scan(&products).Error
	if err != nil {
		return []models.Product{}, err
	}

	return products, nil
}

// func (r *productRepository) CheckStock(product_id uint) (int, error) {
// 	var stock int
// 	if err := r.DB.Raw("SELECT stock FROM products WHERE id=$1", product_id).Scan(&stock).Error; err != nil {
// 		return 0, err
// 	}
// 	return stock, nil
// }

func (r *productRepository) SearchProducts(key string) ([]models.Product, error) {

	var productDetails []models.Product

	query := `SELECT i.* FROM products i LEFT JOIN categories c ON i.category_id = c.id 
          WHERE i.name ILIKE '%' || ? || '%' OR c.name ILIKE '%' || ? || '%'`
	if err := r.DB.Raw(query, key, key).Scan(&productDetails).Error; err != nil {
		return []models.Product{}, err
	}

	return productDetails, nil
}

// func (r *productRepository) SaveSearchHistory(user_id uint, key string) error {

// 	err := r.DB.Exec(`INSERT INTO search_histories (user_id,search_key) VALUES (?,?)`,
// 		user_id, key).Error
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *productRepository) GetSearchHistory(user_id uint) ([]models.SearchHistory, error) {

// 	var searchHistory []models.SearchHistory

// 	err := r.DB.Raw(`SELECT * FROM search_histories WHERE user_id=$1 ORDER BY created_at DESC`,
// 		user_id).Scan(&searchHistory).Error
// 	if err != nil {
// 		return []models.SearchHistory{}, err
// 	}

// 	return searchHistory, nil
// }

func (r *productRepository) UpdateProductImage(product_id uint, url string) (models.Product, error) {

	var updateProduct models.Product

	// Use GORM's Update method to update the product image
	result := r.DB.Model(&domain.Product{}).Where("id = ?", product_id).Update("image", url).Scan(&updateProduct)

	// Check for errors in the update process
	if result.Error != nil {
		return models.Product{}, result.Error
	}

	// Check if any row was affected (to handle case where no rows matched)
	if result.RowsAffected == 0 {
		return models.Product{}, errors.ErrEntityNotFound
	}

	return updateProduct, nil
}

func (r *productRepository) UpdateProduct(product_id uint, model models.Product) (models.Product, error) {

	var updateProduct models.Product

	result := r.DB.Model(&domain.Product{}).Where("id = ?", product_id).Updates(
		domain.Product{
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
		return models.Product{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.Product{}, errors.ErrEntityNotFound
	}

	return updateProduct, nil
}
