package http

import (
	eerrors "errors"
	"reflect"
	"testing"

	nhttp "net/http"
	"net/http/httptest"

	"github.com/go-playground/assert/v2"
	"github.com/rghiorghisor/basic-go-rest-api/errors"
	"github.com/stretchr/testify/mock"

	"github.com/gin-gonic/gin"
)

func TestHandle404(t *testing.T) {

	err := errors.NewEntityNotFound(reflect.TypeOf(&TestStruct{}), "1")
	router, mock := setup(err)

	mock.On("Operation").Return(err)

	w := httptest.NewRecorder()
	req, _ := nhttp.NewRequest("GET", "/api", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestHandle500(t *testing.T) {
	err := eerrors.New("unexpected")
	router, mock := setup(err)

	mock.On("Operation").Return(err)

	w := httptest.NewRecorder()
	req, _ := nhttp.NewRequest("GET", "/api", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
}

func setup(err error) (r *gin.Engine, mock *ErrorsControllerTestMock) {
	router := gin.Default()
	router.Use(
		JSONAppErrorHandler(),
	)

	c := NewController()

	api := router.Group("/api")

	api.GET("", c.Operation)

	return router, c.s
}

type TestStruct struct {
}

type TestController struct {
	s *ErrorsControllerTestMock
}

func NewController() *TestController {
	return &TestController{new(ErrorsControllerTestMock)}
}

func (c *TestController) Operation(ctx *gin.Context) {
	err := c.s.Operation(ctx)

	ctx.Error(err)
}

type ErrorsControllerTestMock struct {
	mock.Mock
}

func (m *ErrorsControllerTestMock) Operation(ctx *gin.Context) error {
	args := m.Called()

	return args.Error(0)
}
