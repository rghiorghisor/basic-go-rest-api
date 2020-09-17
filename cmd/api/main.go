package main

import (
	"github.com/rghiorghisor/basic-go-rest-api/appserver"
	"github.com/rghiorghisor/basic-go-rest-api/config"
	"github.com/rghiorghisor/basic-go-rest-api/container"
	"github.com/rghiorghisor/basic-go-rest-api/logger"

	property_controller "github.com/rghiorghisor/basic-go-rest-api/property/gateway/http"
	property_service "github.com/rghiorghisor/basic-go-rest-api/property/service"
	propertyset_controller "github.com/rghiorghisor/basic-go-rest-api/propertyset/gateway/http"
	propertyset_service "github.com/rghiorghisor/basic-go-rest-api/propertyset/service"
	"github.com/rghiorghisor/basic-go-rest-api/server/http"
	server_storage "github.com/rghiorghisor/basic-go-rest-api/server/storage"
)

func main() {
	appServer := appserver.New()
	appServer.LoadConfig()

	setupStorage(appServer.Configuration, appServer.Container)
	setupServices(appServer.Container)
	setupControllers(appServer.Container)
	setupServer(appServer.Container)

	appServer.Start()
}

func setupStorage(appConfiguration *config.AppConfiguration, c *container.Container) {
	c.Provide(func() (*server_storage.Storage, error) {
		storage := server_storage.New()
		if err := storage.SetupStorage(appConfiguration.Storage); err != nil {
			logger.Main.Error("Cannot setup storage", err)
		}

		return storage, nil
	})
}

func setupServer(c *container.Container) {
	c.Provide(http.NewServer)
}

func setupServices(c *container.Container) {
	c.Provide(property_service.New)
	c.Provide(propertyset_service.New)

	// Add here additional services...
}

func setupControllers(c *container.Container) {
	c.Provide(property_controller.New)
	c.Provide(propertyset_controller.New)

	// Add here additional controllers...
}
