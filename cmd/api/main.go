package main

import (
	"log"

	"ahava/cmd/api/docs"
	"ahava/pkg/config"
	di "ahava/pkg/di"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	// Load the configuration
	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	// Swagger configuration
	docs.SwaggerInfo.Title = "Ahava"
	docs.SwaggerInfo.Description = "An online store for purchasing high-quality jerseys of your favorite clubs."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = config.BASE_URL
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Initialize DI and start the server
	server, diErr := di.InitializeAPI(config)
	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	} else {
		server.Start()
	}
}
