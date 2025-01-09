// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"ahava/pkg/api"
	"ahava/pkg/api/handler"
	"ahava/pkg/config"
	"ahava/pkg/db"
	"ahava/pkg/helper"
	"ahava/pkg/repository"
	"ahava/pkg/service"
)

// Injectors from wire.go:

// InitializeAPI is the entry point for the wire dependency injection setup
func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	gormDB, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}
	userRepository := repository.NewUserRepository(gormDB)
	helperHelper := helper.NewHelper(cfg)
	userService := service.NewUserService(userRepository, cfg, helperHelper)
	userHandler := handler.NewUserHandler(userService)
	adminRepository := repository.NewAdminRepository(gormDB)
	adminService := service.NewAdminService(adminRepository, helperHelper)
	adminHandler := handler.NewAdminHandler(adminService)
	categoryRepository := repository.NewCategoryRepository(gormDB)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productRepository := repository.NewProductRepository(gormDB)
	offerRepository := repository.NewOfferRepository(gormDB)
	productService := service.NewProductService(productRepository, offerRepository, helperHelper)
	productHandler := handler.NewProductHandler(productService)
	orderRepository := repository.NewOrderRepository(gormDB)
	cartRepository := repository.NewCartRepository(gormDB)
	cartService := service.NewCartService(cartRepository, userRepository, offerRepository)
	orderService := service.NewOrderService(orderRepository, cartService)
	orderHandler := handler.NewOrderHandler(orderService)
	cartHandler := handler.NewCartHandler(cartService)
	paymentRepository := repository.NewPaymentRepository(gormDB)
	paymentService := service.NewPaymentService(paymentRepository, orderRepository)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	offerService := service.NewOfferService(offerRepository)
	offerHandler := handler.NewOfferHandler(offerService)
	wishlistRepository := repository.NewWishlistRepository(gormDB)
	wishlistService := service.NewWishlistService(wishlistRepository, offerRepository)
	wishlistHandler := handler.NewWishlistHandler(wishlistService)
	serverHTTP := http.NewServerHTTP(userHandler, adminHandler, categoryHandler, productHandler, orderHandler, cartHandler, paymentHandler, offerHandler, wishlistHandler)
	return serverHTTP, nil
}
