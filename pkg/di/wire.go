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
		// repository.NewOfferRepository,
		repository.NewWishlistRepository,
		repository.NewCartRepository,
		repository.NewOrderRepository,
		repository.NewPaymentRepository,
		repository.NewNewsRepository,

		service.NewUserService,
		service.NewAdminService,
		service.NewProductService,
		// service.NewOfferService,
		service.NewWishlistService,
		service.NewCartService,
		service.NewOrderService,
		service.NewPaymentService,
		service.NewUploadService,
		service.NewNewsService,

		handler.NewUserHandler,
		handler.NewAdminHandler,
		handler.NewProductHandler,
		// handler.NewOfferHandler,
		handler.NewWishlistHandler,
		handler.NewCartHandler,
		handler.NewOrderHandler,
		handler.NewPaymentHandler,
		handler.NewUploadHandler,
		handler.NewNewsHandler,

		helper.NewHelper,

		http.NewServerHTTP,
	)

	return &http.ServerHTTP{}, nil
}
