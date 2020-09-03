package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rghiorghisor/basic-go-rest-api/property/gateway/http"
)

// Controller defines the default functionality of a handler. It must be
// implemented by any struct that needs to be registered further.
type Controller interface {
	Register(routerGroup *gin.RouterGroup)
}

// Controllers is a collections of all available controllers.
type Controllers struct {
	HTTP []Controller
}

// NewControllers creates a new collection of controllers based on the given services.
func NewControllers(services *Services) *Controllers {
	instance := &Controllers{}

	instance.HTTP = append(instance.HTTP, http.NewController(services.PropertyService))

	return instance
}
