package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	// AUTHTOKEN               string `mapstructure:"DB_AUTHTOKEN"`
	// ACCOUNTSID              string `mapstructure:"DB_ACCOUNTSID"`
	// SERVICESID              string `mapstructure:"DB_SERVICESID"`
	MINIO_ENDPOINT          string `mapstructure:"MINIO_ENDPOINT"`
	MINIO_ENDPOINT_PUBLIC   string `mapstructure:"MINIO_ENDPOINT_PUBLIC"`
	MINIO_ACCESS_KEY_ID     string `mapstructure:"MINIO_ACCESS_KEY_ID"`
	MINIO_SECRET_ACCESS_KEY string `mapstructure:"MINIO_SECRET_ACCESS_KEY"`
}

var envs = []string{
	"DB_HOST",
	"DB_NAME",
	"DB_USER",
	"DB_PORT",
	"DB_PASSWORD",
	// "DB_AUTHTOKEN",
	// "DB_ACCOUNTSID",
	// "DB_SERVICESID",
	"MINIO_ENDPOINT",
	"MINIO_ENDPOINT_PUBLIC",
	"MINIO_ACCESS_KEY_ID",
	"MINIO_SECRET_ACCESS_KEY",
}

func LoadConfig() (Config, error) {
	var config Config

	viper.AutomaticEnv()

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if err := validator.New().Struct(&config); err != nil {
		return config, err
	}

	fmt.Println(config)

	return config, nil
}
