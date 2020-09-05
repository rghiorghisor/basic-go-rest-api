package http

import (
	net_http "net/http"

	"github.com/gin-gonic/gin"
)

// HealthcheckController handles any such check operations.
type HealthcheckController struct {
}

// NewHealthcheckController retrieves a new controller that handles any health check requests.
func NewHealthcheckController() *HealthcheckController {
	return new(HealthcheckController)
}

func (ctrl *HealthcheckController) check(ctx *gin.Context) {
	ctx.String(net_http.StatusOK, "OK")
}

// Register this controller to the provided group.
func (ctrl *HealthcheckController) Register(routerGroup *gin.RouterGroup) {

	routerGroup.GET("/healthcheck", ctrl.check)

}
