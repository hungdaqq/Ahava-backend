//go:build wireinject
// +build wireinject

package di

import (
	http "ahava/pkg/api"
	handler "ahava/pkg/api/handler"
	config "ahava/pkg/config"
	db "ahava/pkg/db"
	repository "ahava/pkg/repository"
	usecase "ahava/pkg/usecase"
	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(db.ConnectDatabase, repository.NewUserRepository, usecase.NewUserUseCase, handler.NewUserHandler, http.NewServerHTTP)

	return &http.ServerHTTP{}, nil
}
