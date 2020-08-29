package main

import (
	"github.com/rghiorghisor/basic-go-rest-api/config"
	"github.com/rghiorghisor/basic-go-rest-api/logger"
	"github.com/rghiorghisor/basic-go-rest-api/server"
	"github.com/rghiorghisor/basic-go-rest-api/server/storage"
)

func main() {
	// Load and validate configuration.
	appConfiguration := config.NewAppConfiguration()
	appConfiguration.Load()

	// Setup and start the logger.
	logger.New(appConfiguration.Logger)
	logger.Main.Info("Starting application...")

	// Configure and connect to storage.
	storage := storage.NewStorage()
	storage.SetupStorage(appConfiguration.Storage.DbConfiguration)

	// Configure and setup services.
	services := server.NewServices()
	services.SetupServices(storage)

	// Create and run the server.
	appServer := server.NewAppServer()
	appServer.Setup(appConfiguration.Server.HTTPServer, services)
	appServer.Run()
}
