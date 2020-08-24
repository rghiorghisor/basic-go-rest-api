package server

import (
	"github.com/rghiorghisor/basic-go-rest-api/property"
	pservice "github.com/rghiorghisor/basic-go-rest-api/property/service"
	"github.com/rghiorghisor/basic-go-rest-api/server/storage"
)

// Services structure holds all available services.
type Services struct {
	PropertyService property.Service
}

// NewServices retrieves a new Services struct.
func NewServices() *Services {
	return &Services{}
}

// SetupServices creates and initializes all available services.
func (services *Services) SetupServices(storage *storage.Storage) {
	services.PropertyService = pservice.NewPropertyService(storage)
}
