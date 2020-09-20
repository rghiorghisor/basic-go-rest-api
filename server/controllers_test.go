package server

import (
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

func TestControllers(t *testing.T) {
	controllerParams := ControllersParams{}
	controllerParams.Controllers = append(controllerParams.Controllers, newTestController1().Controller)
	controllerParams.Controllers = append(controllerParams.Controllers, newTestController2().Controller)

	controllers := NewControllersWithParams(controllerParams)

	assert.Equal(t, 2, len(controllers.HTTP))
}

type testController1 struct {
}

func newTestController1() ControllerWrapper {
	return ControllerWrapper{
		Controller: &testController1{},
	}
}

func (t *testController1) Register(routerGroup *gin.RouterGroup) {

}

type testController2 struct {
}

func newTestController2() ControllerWrapper {
	return ControllerWrapper{
		Controller: &testController2{},
	}
}

func (t *testController2) Register(routerGroup *gin.RouterGroup) {

}
