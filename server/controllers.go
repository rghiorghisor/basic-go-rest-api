package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
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

// NewControllersWithParams creates a new collection of controllers based on the given services.
func NewControllersWithParams(cp ControllersParams) Controllers {
	instance := Controllers{}

	for _, controller := range cp.Controllers {
		instance.HTTP = append(instance.HTTP, controller)
	}

	return instance
}

// ControllerWrapper contains a Controller along and helps to initialize the Controllers.
type ControllerWrapper struct {
	dig.Out

	Controller Controller `group:"controllers"`
}

// ControllersParams contains all managed controllers.
type ControllersParams struct {
	dig.In

	Controllers []Controller `group:"controllers"`
}
