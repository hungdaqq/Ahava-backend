package repository

import (
	"ahava/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type ProductRepository interface {
	AddProduct(product models.Products) (models.Products, error)

	UpdateProduct(int, models.Products) (models.Products, error)
	UpdateProductImage(int, string) (models.Products, error)

	DeleteProduct(product_id string) error

	ShowProductDetails(product_id int) (models.Products, error)
	ListProducts(limit, offset int) (models.ListProducts, error)
	ListCategoryProducts() ([]models.CategoryProducts, error)

	// ListProductsByCategory(id int) ([]models.Products, error)
	CheckStock(product_id int) (int, error)
	// CheckPrice(product_id int) (float64, error)
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

func (i *productRepository) CheckProduct(pid int) (bool, error) {
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
func (i *productRepository) ShowProductDetails(product_id int) (models.Products, error) {

	var product models.Products

	err := i.DB.Raw(`SELECT * FROM products WHERE products.id = ?`,
		product_id).Scan(&product).Error
	if err != nil {
		return models.Products{}, errors.New("no records with that ID exist")
	}

	return product, nil
}

func (ad *productRepository) ListProducts(limit, offset int) (models.ListProducts, error) {

	var listProducts models.ListProducts
	var productDetails []models.Products
	var total int64

	query := ad.DB.Model(&models.Products{})
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

type ProductWithCategory struct {
	models.Products
	CategoryName string `json:"category_name"`
}

func (ad *productRepository) ListCategoryProducts() ([]models.CategoryProducts, error) {

	var productWithCategories []ProductWithCategory
	err := ad.DB.Raw("SELECT products.*, categories.category_name FROM products JOIN categories ON categories.id = products.category_id").
		Scan(&productWithCategories).Error
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[int]*models.CategoryProducts)
	for _, product := range productWithCategories {
		if _, exists := categoryMap[product.CategoryID]; !exists {
			categoryMap[product.CategoryID] = &models.CategoryProducts{
				CategoryID:   product.CategoryID,
				CategoryName: product.CategoryName,
				Products:     []models.Products{},
			}
		}
		categoryMap[product.CategoryID].Products = append(categoryMap[product.CategoryID].Products, product.Products)
	}

	var categoryProducts []models.CategoryProducts
	for _, category := range categoryMap {
		categoryProducts = append(categoryProducts, *category)
	}
	return categoryProducts, nil
}

func (i *productRepository) CheckStock(product_id int) (int, error) {
	var stock int
	if err := i.DB.Raw("SELECT stock FROM products WHERE id=$1", product_id).Scan(&stock).Error; err != nil {
		return 0, err
	}
	return stock, nil
}

func (i *productRepository) CheckPrice(pid int) (float64, error) {
	var k float64
	err := i.DB.Raw("SELECT price FROM products WHERE id=?", pid).Scan(&k).Error
	if err != nil {
		return 0, err
	}

	return k, nil
}

func (ad *productRepository) SearchProducts(key string) ([]models.Products, error) {

	var productDetails []models.Products

	query := `SELECT i.* FROM products i LEFT JOIN categories c ON i.category_id = c.id 
				WHERE i.product_name ILIKE '%' || ? || '%' OR c.category_name ILIKE '%' || ? || '%'`
	if err := ad.DB.Raw(query, key, key).Scan(&productDetails).Error; err != nil {
		return []models.Products{}, err
	}

	return productDetails, nil
}

func (i *productRepository) UpdateProductImage(product_id int, url string) (models.Products, error) {

	var updateProduct models.Products

	result := i.DB.Raw("UPDATE products SET image = $1 WHERE id = $2 RETURNING *",
		url, product_id).Scan(&updateProduct)
	if result.Error != nil {
		return models.Products{}, result.Error
	}
	if result.RowsAffected < 1 {
		return models.Products{}, errors.New("no records with that ID exist")
	}

	return updateProduct, nil
}

func (i *productRepository) UpdateProduct(product_id int, model models.Products) (models.Products, error) {

	var updateProduct models.Products

	result := i.DB.Raw("UPDATE products SET product_name=$1,category_id=$2,price=$3,size=$4,stock=$5,description=$6,how_to_use=$7 WHERE id = $8 RETURNING *",
		model.ProductName, model.CategoryID, model.Price, model.Size, model.Stock, model.Description, model.HowToUse, product_id).Scan(&updateProduct)
	if result.Error != nil {
		return models.Products{}, result.Error
	}
	if result.RowsAffected < 1 {
		return models.Products{}, errors.New("no records with that ID exist")
	}

	return updateProduct, nil
}
