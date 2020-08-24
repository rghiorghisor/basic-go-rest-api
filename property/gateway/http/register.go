package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rghiorghisor/basic-go-rest-api/property"
)

// RegisterEndpoints makes sure the available operations are exposed.
func RegisterEndpoints(routerGroup *gin.RouterGroup, service property.Service) {
	controller := NewController(service)

	api := routerGroup.Group("/property")

	api.POST("", controller.Create)
	api.GET("", controller.ReadAll)
	api.GET("/:id", controller.Read)
	api.PUT("/:id", controller.Update)
	api.DELETE("/:id", controller.Delete)
}
