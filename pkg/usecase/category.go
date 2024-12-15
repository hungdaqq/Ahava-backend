package usecase

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
	"fmt"
)

type CategoryUseCase interface {
	AddCategory(category models.Category) (models.Category, error)
	UpdateCategory(categoryID int, category, description string) (models.Category, error)
	DeleteCategory(categoryID int) error
	GetCategories() ([]models.Category, error)
	GetBannersForUsers() ([]models.Banner, error)
}

type categoryUseCase struct {
	repository repository.CategoryRepository
	// productRepository repository.ProductRepository
	// offerRepository   repository.OfferRepository
}

func NewCategoryUseCase(
	repo repository.CategoryRepository,
	// inv interfaces.ProductRepository,
	// offer interfaces.OfferRepository,
) CategoryUseCase {
	return &categoryUseCase{
		repository: repo,
		// productRepository: inv,
		// offerRepository:   offer,
	}
}

func (Cat *categoryUseCase) AddCategory(category models.Category) (models.Category, error) {

	result, err := Cat.repository.AddCategory(category)

	if err != nil {
		return models.Category{}, err
	}

	return result, nil

}
func (Cat *categoryUseCase) GetCategories() ([]models.Category, error) {

	result, err := Cat.repository.GetCategories()

	if err != nil {
		return []models.Category{}, err
	}

	return result, nil

}

func (Cat *categoryUseCase) UpdateCategory(categoryID int, category, description string) (models.Category, error) {

	newcat, err := Cat.repository.UpdateCategory(categoryID, category, description)
	if err != nil {
		return models.Category{}, err
	}

	return newcat, err
}

func (Cat *categoryUseCase) DeleteCategory(categoryID int) error {

	err := Cat.repository.DeleteCategory(categoryID)
	if err != nil {
		return err
	}
	return nil

}

func (Cat *categoryUseCase) GetBannersForUsers() ([]models.Banner, error) {
	// Find categories with the highest offer percentage, at least one, maximum 3.
	banners, err := Cat.repository.GetBannersForUsers()
	if err != nil {
		return nil, err
	}

	// Find images of 2 products from each category.
	for i := range banners {
		images, err := Cat.repository.GetImagesOfProductsFromACategory(banners[i].CategoryID)
		if err != nil {
			return nil, err
		}
		banners[i].Images = images
		fmt.Println("loop instance", banners[i])
	}

	fmt.Println("banners", banners)
	return banners, nil
}
