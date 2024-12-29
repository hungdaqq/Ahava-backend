package usecase

import (
	helper "ahava/pkg/helper"
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
	"mime/multipart"
)

type ProductUseCase interface {
	AddProduct(product models.Products, default_image *multipart.FileHeader, images []*multipart.FileHeader) (models.Products, error)
	UpdateProduct(uint, models.Products) (models.Products, error)
	UpdateProductImage(product_id uint, file *multipart.FileHeader) (models.Products, error)
	DeleteProduct(product_id uint) error
	GetProductDetails(product_id uint) (models.Products, error)
	ListCategoryProducts(category_id uint) ([]models.Products, error)
	ListProductsForAdmin(limit, offest int) (models.ListProducts, error)
	ListFeaturedProducts() ([]models.Products, error)

	SearchProducts(key string) ([]models.Products, error)
	// GetSearchHistory(user_id uint) ([]models.SearchHistory, error)
}

type productUseCase struct {
	repository      repository.ProductRepository
	offerRepository repository.OfferRepository
	helper          helper.Helper
	// wishlistRepository repository.WishlistRepository
}

func NewProductUseCase(
	repo repository.ProductRepository,
	offer repository.OfferRepository,
	h helper.Helper,
	// w repository.WishlistRepository,
) *productUseCase {
	return &productUseCase{
		repository:      repo,
		offerRepository: offer,
		helper:          h,
		// wishlistRepository: w,
	}
}

func (i *productUseCase) AddProduct(product models.Products, default_image *multipart.FileHeader, images []*multipart.FileHeader) (models.Products, error) {

	default_image_url, err := i.helper.AddImageToS3(default_image, "ahava")
	if err != nil {
		return models.Products{}, err
	}

	var urls string

	for idx, image := range images {
		url, err := i.helper.AddImageToS3(image, "ahava")
		if err != nil {
			return models.Products{}, err
		}

		if idx > 0 {
			urls += "," // Add a comma before every subsequent URL (except the first one)
		}
		urls += url // Append the URL
	}

	product.DefaultImage = default_image_url
	product.Images = urls
	result, err := i.repository.AddProduct(product)
	if err != nil {
		return models.Products{}, err
	}

	return result, nil
}

func (i *productUseCase) UpdateProductImage(id uint, file *multipart.FileHeader) (models.Products, error) {

	url, err := i.helper.AddImageToS3(file, "ahava")
	if err != nil {
		return models.Products{}, err
	}

	//send the url and save it in database
	result, err := i.repository.UpdateProductImage(id, url)
	if err != nil {
		return models.Products{}, err
	}

	return result, nil

}

func (i *productUseCase) UpdateProduct(id uint, model models.Products) (models.Products, error) {

	//send the url and save it in database
	result, err := i.repository.UpdateProduct(id, model)
	if err != nil {
		return models.Products{}, err
	}

	return result, nil
}

func (i *productUseCase) DeleteProduct(product_id uint) error {

	err := i.repository.DeleteProduct(product_id)
	if err != nil {
		return err
	}
	return nil
}

func (i *productUseCase) GetProductDetails(product_id uint) (models.Products, error) {

	product, err := i.repository.GetProductDetails(product_id)
	if err != nil {
		return models.Products{}, err
	}

	offerPercentage, err := i.offerRepository.FindOfferRate(product.ID)
	if err != nil {
		return models.Products{}, err
	}

	if offerPercentage > 0 {
		product.DiscountedPrice = product.Price - (product.Price*uint64(offerPercentage))/100
	} else {
		product.DiscountedPrice = product.Price
	}

	return product, nil
}

func (i *productUseCase) ListCategoryProducts(category_id uint) ([]models.Products, error) {

	products, err := i.repository.ListCategoryProducts(category_id)
	if err != nil {
		return nil, err
	}

	for idx := range products {
		offerPercentage, err := i.offerRepository.FindOfferRate(products[idx].CategoryID)
		if err != nil {
			return nil, err
		}
		if offerPercentage > 0 {
			products[idx].DiscountedPrice = products[idx].Price - (products[idx].Price*uint64(offerPercentage))/100
		} else {
			products[idx].DiscountedPrice = products[idx].Price
		}
	}

	return products, nil
}

func (i *productUseCase) ListFeaturedProducts() ([]models.Products, error) {

	products, err := i.repository.ListFeaturedProducts()
	if err != nil {
		return []models.Products{}, err
	}

	for idx := range products {
		offerPercentage, err := i.offerRepository.FindOfferRate(products[idx].CategoryID)
		if err != nil {
			return nil, err
		}
		if offerPercentage > 0 {
			products[idx].DiscountedPrice = products[idx].Price - (products[idx].Price*uint64(offerPercentage))/100
		} else {
			products[idx].DiscountedPrice = products[idx].Price
		}
	}

	return products, nil
}

// func (i *productUseCase) ListProductsForUser(page, userID uint) ([]models.Products, error) {

// 	productDetails, err := i.repository.ListProducts(page)
// 	if err != nil {
// 		return []models.Products{}, err
// 	}

// 	fmt.Println("product details is:", productDetails)

// 	//loop inside products and then calculate discounted price of each then return
// 	for j := range productDetails {
// 		discount_percentage, err := i.offerRepository.FindDiscountPercentage(productDetails[j].CategoryID)
// 		if err != nil {
// 			return []models.Products{}, errors.New("there was some error in finding the discounted prices")
// 		}
// 		var discount float64

// 		if discount_percentage > 0 {
// 			discount = (productDetails[j].Price * float64(discount_percentage)) / 100
// 		}

// 		productDetails[j].DiscountedPrice = productDetails[j].Price - discount

// 		productDetails[j].IfPresentAtWishlist, err = i.wishlistRepository.CheckIfTheItemIsPresentAtWishlist(userID, int(productDetails[j].ID))
// 		if err != nil {
// 			return []models.Products{}, errors.New("error while checking ")
// 		}

// 		productDetails[j].IfPresentAtCart, err = i.wishlistRepository.CheckIfTheItemIsPresentAtCart(userID, int(productDetails[j].ID))
// 		if err != nil {
// 			return []models.Products{}, errors.New("error while checking ")
// 		}

// 	}

// 	return productDetails, nil
// }

func (i *productUseCase) ListProductsForAdmin(limit, offset int) (models.ListProducts, error) {

	listProducts, err := i.repository.ListProducts(limit, offset)
	if err != nil {
		return models.ListProducts{}, err
	}

	for idx := range listProducts.Products {
		offerPercentage, err := i.offerRepository.FindOfferRate(listProducts.Products[idx].CategoryID)
		if err != nil {
			return models.ListProducts{}, err
		}
		if offerPercentage > 0 {
			listProducts.Products[idx].DiscountedPrice = listProducts.Products[idx].Price - (listProducts.Products[idx].Price*uint64(offerPercentage))/100
		} else {
			listProducts.Products[idx].DiscountedPrice = listProducts.Products[idx].Price
		}
	}

	return listProducts, nil
}

func (i *productUseCase) SearchProducts(key string) ([]models.Products, error) {

	products, err := i.repository.SearchProducts(key)
	if err != nil {
		return []models.Products{}, err
	}

	for idx := range products {
		offerPercentage, err := i.offerRepository.FindOfferRate(products[idx].CategoryID)
		if err != nil {
			return nil, err
		}
		if offerPercentage > 0 {
			products[idx].DiscountedPrice = products[idx].Price - (products[idx].Price*uint64(offerPercentage))/100
		} else {
			products[idx].DiscountedPrice = products[idx].Price
		}
	}

	return products, nil
}

// func (i *productUseCase) GetSearchHistory(user_id uint) ([]models.SearchHistory, error) {

// 	searchHistory, err := i.repository.GetSearchHistory(user_id)
// 	if err != nil {
// 		return []models.SearchHistory{}, err
// 	}

// 	return searchHistory, nil

// }
