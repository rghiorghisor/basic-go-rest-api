package main

import (
	"log"

	"github.com/rghiorghisor/basic-go-rest-api/config"
	"github.com/rghiorghisor/basic-go-rest-api/logger"
	"github.com/rghiorghisor/basic-go-rest-api/server"
	"github.com/rghiorghisor/basic-go-rest-api/server/http"
	"github.com/rghiorghisor/basic-go-rest-api/server/storage"
)

func main() {
	// Load and validate configuration.
	appConfiguration := config.NewAppConfiguration()
	appConfiguration.Load()

	// Setup and start the logger.
	startLogger(appConfiguration)

	// Configure and connect to storage.
	storage := storage.NewStorage()
	storage.SetupStorage(appConfiguration.Storage.DbConfiguration)

	// Configure and setup services.
	services := server.NewServices(storage)

	// Gather the available controllers.
	controllers := server.NewControllers(services)

	// Create and run the server.
	appServer := http.NewAppServer()
	appServer.Setup(appConfiguration.Server.HTTPServer, controllers)

	if err := appServer.Run(); err != nil {
		log.Fatalf("Failed to start: %+v", err)
	}
}

func startLogger(appConfiguration *config.AppConfiguration) {
	logger.New(appConfiguration.Logger)

	logger.Main.Info("Loaded configuration.")
	logger.Main.Infof("Starting application in %s mode...", appConfiguration.Environment.Name)
}
