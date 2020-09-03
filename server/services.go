package server

import (
	"github.com/rghiorghisor/basic-go-rest-api/property"
	property_service "github.com/rghiorghisor/basic-go-rest-api/property/service"
	"github.com/rghiorghisor/basic-go-rest-api/server/storage"
)

// Services structure holds all available services.
type Services struct {
	PropertyService property.Service
	// Define here new services.
}

// NewServices retrieves a new Services struct.
func NewServices(storage *storage.Storage) *Services {
	return &Services{
		property_service.NewPropertyService(storage),
		// Initialize here new services.
	}
}
