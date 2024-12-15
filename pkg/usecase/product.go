package usecase

import (
	helper "ahava/pkg/helper"
	repository "ahava/pkg/repository"
	"ahava/pkg/utils/models"
	"mime/multipart"
)

type ProductUseCase interface {
	AddProduct(product models.Products, image *multipart.FileHeader) (models.Products, error)
	UpdateProduct(int, models.Products) (models.Products, error)
	UpdateProductImage(product_id int, file *multipart.FileHeader) (models.Products, error)
	DeleteProduct(product_id string) error

	ShowProductDetails(product_id int) (models.Products, error)
	ListProductsForUser() ([]models.CategoryProducts, error)
	ListProductsForAdmin(limit, offest int) (models.ListProducts, error)

	SearchProducts(key string) ([]models.Products, error)
}

type productUseCase struct {
	repository repository.ProductRepository
	// offerRepository    repository.OfferRepository
	helper helper.Helper
	// wishlistRepository repository.WishlistRepository
}

func NewProductUseCase(
	repo repository.ProductRepository,
	// offer repository.OfferRepository,
	h helper.Helper,
	// w repository.WishlistRepository,
) *productUseCase {
	return &productUseCase{
		repository: repo,
		// offerRepository:    offer,
		helper: h,
		// wishlistRepository: w,
	}
}

func (i *productUseCase) AddProduct(product models.Products, image *multipart.FileHeader) (models.Products, error) {

	url, err := i.helper.AddImageToS3(image)
	if err != nil {
		return models.Products{}, err
	}

	product.Image = url
	ProductResponse, err := i.repository.AddProduct(product)
	if err != nil {
		return models.Products{}, err
	}

	return ProductResponse, nil

}

func (i *productUseCase) UpdateProductImage(id int, file *multipart.FileHeader) (models.Products, error) {

	url, err := i.helper.AddImageToS3(file)
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

func (i *productUseCase) UpdateProduct(id int, model models.Products) (models.Products, error) {

	//send the url and save it in database
	result, err := i.repository.UpdateProduct(id, model)
	if err != nil {
		return models.Products{}, err
	}

	return result, nil

}

func (i *productUseCase) DeleteProduct(productID string) error {

	err := i.repository.DeleteProduct(productID)
	if err != nil {
		return err
	}
	return nil

}

func (i *productUseCase) ShowProductDetails(id int) (models.Products, error) {

	product, err := i.repository.ShowProductDetails(id)
	if err != nil {
		return models.Products{}, err
	}

	// DiscountPercentage, err := i.offerRepository.FindDiscountPercentage(product.CategoryID)
	// if err != nil {
	// 	return models.Products{}, err
	// }

	// //make discounted price by calculation
	// var discount float64
	// if DiscountPercentage > 0 {
	// 	discount = (product.Price * float64(DiscountPercentage)) / 100
	// }

	// product.DiscountedPrice = product.Price - discount

	return product, nil

}

func (i *productUseCase) ListProductsForUser() ([]models.CategoryProducts, error) {

	categoryProducts, err := i.repository.ListCategoryProducts()
	if err != nil {
		return []models.CategoryProducts{}, err
	}

	return categoryProducts, nil
}

// func (i *productUseCase) ListProductsForUser(page, userID int) ([]models.Products, error) {

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

	// //loop inside products and then calculate discounted price of each then return
	// for j := range listProducts.Products {
	// 	discount_percentage, err := i.offerRepository.FindDiscountPercentage(listProducts.Products[j].CategoryID)
	// 	if err != nil {
	// 		return models.ListProducts{}, errors.New("there was some error in finding the discounted prices")
	// 	}
	// 	var discount float64

	// 	if discount_percentage > 0 {
	// 		discount = (listProducts.Products[j].Price * float64(discount_percentage)) / 100
	// 	}

	// 	listProducts.Products[j].DiscountedPrice = listProducts.Products[j].Price - discount

	// }

	return listProducts, nil

}

func (i *productUseCase) SearchProducts(key string) ([]models.Products, error) {

	productDetails, err := i.repository.SearchProducts(key)
	if err != nil {
		return []models.Products{}, err
	}

	// //loop inside products and then calculate discounted price of each then return
	// for j := range productDetails {
	// 	discount_percentage, err := i.offerRepository.FindDiscountPercentage(productDetails[j].CategoryID)
	// 	if err != nil {
	// 		return []models.Products{}, errors.New("there was some error in finding the discounted prices")
	// 	}
	// 	var discount float64

	// 	if discount_percentage > 0 {
	// 		discount = (productDetails[j].Price * float64(discount_percentage)) / 100
	// 	}

	// 	productDetails[j].DiscountedPrice = productDetails[j].Price - discount

	// }

	return productDetails, nil

}
