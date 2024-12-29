package usecase

import (
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
)

type CategoryUseCase interface {
	AddCategory(category models.Category) (models.Category, error)
	UpdateCategory(category_id uint, category, description string) (models.Category, error)
	DeleteCategory(category_id uint) error
	GetCategories() ([]models.Category, error)
	// GetBannersForUsers() ([]models.Banner, error)
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

func (Cat *categoryUseCase) UpdateCategory(category_id uint, name, description string) (models.Category, error) {

	result, err := Cat.repository.UpdateCategory(category_id, name, description)
	if err != nil {
		return models.Category{}, err
	}

	return result, err
}

func (Cat *categoryUseCase) DeleteCategory(category_id uint) error {

	err := Cat.repository.DeleteCategory(category_id)
	if err != nil {
		return err
	}
	return nil

}

// func (Cat *categoryUseCase) GetBannersForUsers() ([]models.Banner, error) {
// 	// Find categories with the highest offer percentage, at least one, maximum 3.
// 	banners, err := Cat.repository.GetBannersForUsers()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Find images of 2 products from each category.
// 	for i := range banners {
// 		images, err := Cat.repository.GetImagesOfProductsFromACategory(banners[i].CategoryID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		banners[i].Images = images
// 		fmt.Println("loop instance", banners[i])
// 	}

// 	fmt.Println("banners", banners)
// 	return banners, nil
// }
