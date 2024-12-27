//go:build wireinject
// +build wireinject

package di

import (
	http "ahava/pkg/api"
	"ahava/pkg/api/handler"
	config "ahava/pkg/config"
	db "ahava/pkg/db"

	"github.com/google/wire"
)

// InitializeAPI is the entry point for the wire dependency injection setup
func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(
		db.ConnectDatabase, // Provides *db.Database
		// repository.NewUserRepository, // Provides repository.UserRepository
		// usecase.NewAdminUseCase,      // Other use case providers...
		handler.NewUserHandler, // User handler
		// handler.NewAdminHandler, // Admin handler
		http.NewServerHTTP, // HTTP server setup
	)

	// The return value is automatically generated by Wire
	return &http.ServerHTTP{}, nil
}
