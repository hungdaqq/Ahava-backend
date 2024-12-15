package repository

import (
	"ahava/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	AddCategory(category models.Category) (models.Category, error)
	UpdateCategory(categoryID int, category, description string) (models.Category, error)
	DeleteCategory(categoryID int) error
	GetCategories() ([]models.Category, error)
	GetBannersForUsers() ([]models.Banner, error)
	GetImagesOfProductsFromACategory(CategoryID int) ([]string, error)
}

type categoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(DB *gorm.DB) CategoryRepository {
	return &categoryRepository{DB}
}

func (p *categoryRepository) AddCategory(category models.Category) (models.Category, error) {

	var addCategory models.Category

	if err := p.DB.Raw(`INSERT INTO categories (category_name, description) VALUES (?, ?) RETURNING *`,
		category.CategoryName, category.Description).Scan(&addCategory).Error; err != nil {
		return models.Category{}, err
	}

	return addCategory, nil
}

func (p *categoryRepository) UpdateCategory(categoryID int, category, description string) (models.Category, error) {

	var updateCategory models.Category

	result := p.DB.Raw("UPDATE categories SET category_name = $1, description = $2 WHERE id = $3 RETURNING *",
		category, description, categoryID).Scan(&updateCategory)
	if result.Error != nil {
		return models.Category{}, result.Error
	}
	if result.RowsAffected < 1 {
		return models.Category{}, errors.New("no records with that ID exist")
	}

	return updateCategory, nil
}

func (c *categoryRepository) DeleteCategory(categoryID int) error {

	result := c.DB.Exec("DELETE FROM categories WHERE id = ?", categoryID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return errors.New("no records with that ID exist")
	}

	return nil
}

func (c *categoryRepository) GetCategories() ([]models.Category, error) {
	var model []models.Category
	err := c.DB.Raw("SELECT * FROM categories").Scan(&model).Error
	if err != nil {
		return []models.Category{}, err
	}

	return model, nil
}

func (c *categoryRepository) GetBannersForUsers() ([]models.Banner, error) {
	var banners []models.Banner
	err := c.DB.Raw(`select offers.category_id,categories.category as category_name,offers.discount_rate as discount_percentage
	 from offers
	 join categories on categories.id = offers.category_id
	 where offers.discount_rate > 10 
	 Order by offers.discount_rate desc
	 limit 3`).Scan(&banners).Error
	if err != nil {
		return []models.Banner{}, err
	}
	return banners, nil
}

func (c *categoryRepository) GetImagesOfProductsFromACategory(CategoryID int) ([]string, error) {
	var images []string
	err := c.DB.Raw("select image from products where category_id = $1 limit 2", CategoryID).Scan(&images).Error
	if err != nil {
		return []string{}, err
	}

	return images, nil

}
