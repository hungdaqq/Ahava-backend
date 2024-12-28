package repository

import (
	"ahava/pkg/domain"
	"ahava/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type ProductRepository interface {
	AddProduct(product models.Products) (models.Products, error)

	UpdateProduct(uint, models.Products) (models.Products, error)
	UpdateProductImage(uint, string) (models.Products, error)

	DeleteProduct(product_id string) error

	ShowProductDetails(product_id uint) (models.Products, error)
	ListProducts(limit, offset int) (models.ListProducts, error)
	ListCategoryProducts(category_id uint) (models.Products, error)

	// ListProductsByCategory(id uint) ([]models.Products, error)
	// CheckStock(product_id uint) (uint, error)
	// CheckPrice(product_id uint) (uint64, error)
	SearchProducts(key string) ([]models.Products, error)
	SaveSearchHistory(user_id uint, key string) error
	GetSearchHistory(user_id uint) ([]models.SearchHistory, error)
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

	var productResponse models.Products

	err := i.DB.Raw(`
		INSERT INTO products (category_id, product_name, size, stock, price, image, description, how_to_use)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING *`,
		product.CategoryID, product.ProductName, product.Size, product.Stock, product.Price, product.Image, product.Description, product.HowToUse).
		Scan(&productResponse).Error

	if err != nil {
		return models.Products{}, err
	}

	return productResponse, nil

}

func (i *productRepository) CheckProduct(pid uint) (bool, error) {
	var k int
	err := i.DB.Raw("SELECT COUNT(*) FROM products WHERE id=?",
		pid).Scan(&k).Error
	if err != nil {
		return false, err
	}

	if k == 0 {
		return false, err
	}

	return true, err
}

func (i *productRepository) DeleteProduct(productID string) error {

	result := i.DB.Exec("DELETE FROM products WHERE id = ?", productID)

	if result.RowsAffected < 1 {
		return errors.New("no records with that ID exist")
	}

	return nil
}

// detailed product details
func (i *productRepository) ShowProductDetails(product_id uint) (models.Products, error) {

	var product models.Products

	err := i.DB.Raw(`SELECT * FROM products WHERE products.id = ?`,
		product_id).Scan(&product).Error
	if err != nil {
		return models.Products{}, errors.New("no records with that ID exist")
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

func (i *productRepository) ListCategoryProducts(category_id uint) (models.Products, error) {

	var products models.Products
	err := i.DB.Raw("SELECT * FROM products WHERE category_id=$1", category_id).
		Scan(&products).Error
	if err != nil {
		return models.Products{}, err
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

func (i *productRepository) CheckPrice(pid uint) (uint64, error) {
	var k uint64
	err := i.DB.Raw("SELECT price FROM products WHERE id=?", pid).Scan(&k).Error
	if err != nil {
		return 0, err
	}

	return k, nil
}

func (i *productRepository) SearchProducts(key string) ([]models.Products, error) {

	var productDetails []models.Products

	query := `SELECT i.* FROM products i LEFT JOIN categories c ON i.category_id = c.id 
				WHERE i.product_name ILIKE '%' || ? || '%' OR c.category_name ILIKE '%' || ? || '%'`
	if err := i.DB.Raw(query, key, key).Scan(&productDetails).Error; err != nil {
		return []models.Products{}, err
	}

	return productDetails, nil
}

func (i *productRepository) SaveSearchHistory(user_id uint, key string) error {

	err := i.DB.Exec(`INSERT INTO search_histories (user_id,search_key) VALUES (?,?)`,
		user_id, key).Error
	if err != nil {
		return err
	}

	return nil
}

func (i *productRepository) GetSearchHistory(user_id uint) ([]models.SearchHistory, error) {

	var searchHistory []models.SearchHistory

	err := i.DB.Raw(`SELECT * FROM search_histories WHERE user_id=$1 ORDER BY created_at DESC`,
		user_id).Scan(&searchHistory).Error
	if err != nil {
		return []models.SearchHistory{}, err
	}

	return searchHistory, nil
}

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
		return models.Products{}, errors.New("product not found")
	}

	return updateProduct, nil
}

func (i *productRepository) UpdateProduct(product_id uint, model models.Products) (models.Products, error) {

	var updateProduct models.Products

	result := i.DB.Model(&domain.Products{}).Where("id = ?", product_id).Updates(
		domain.Products{
			ProductName: model.ProductName,
			CategoryID:  model.CategoryID,
			Price:       model.Price,
			Size:        model.Size,
			Stock:       model.Stock,
			Description: model.Description,
			HowToUse:    model.HowToUse,
		}).Scan(&updateProduct)

	if result.Error != nil {
		return models.Products{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.Products{}, errors.New("product not found")
	}

	return updateProduct, nil
}
