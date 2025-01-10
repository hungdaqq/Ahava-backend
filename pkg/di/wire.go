//go:build wireinject
// +build wireinject

package di

import (
	http "ahava/pkg/api"
	"ahava/pkg/api/handler"
	config "ahava/pkg/config"
	db "ahava/pkg/db"
	"ahava/pkg/helper"
	"ahava/pkg/repository"
	"ahava/pkg/service"

	"github.com/google/wire"
)

// InitializeAPI is the entry point for the wire dependency injection setup
func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(
		db.ConnectDatabase,
		repository.NewUserRepository,
		repository.NewAdminRepository,
		repository.NewProductRepository,
		repository.NewCategoryRepository,
		repository.NewOfferRepository,
		repository.NewWishlistRepository,
		repository.NewCartRepository,
		repository.NewOrderRepository,
		repository.NewPaymentRepository,

		service.NewUserService,
		service.NewAdminService,
		service.NewProductService,
		service.NewCategoryService,
		service.NewOfferService,
		service.NewWishlistService,
		service.NewCartService,
		service.NewOrderService,
		service.NewPaymentService,

		handler.NewUserHandler,
		handler.NewAdminHandler,
		handler.NewProductHandler,
		handler.NewCategoryHandler,
		handler.NewOfferHandler,
		handler.NewWishlistHandler,
		handler.NewCartHandler,
		handler.NewOrderHandler,
		handler.NewPaymentHandler,

		helper.NewHelper,

		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil
}
