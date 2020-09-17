package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	apperrors "github.com/rghiorghisor/basic-go-rest-api/errors"
	"github.com/rghiorghisor/basic-go-rest-api/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	router, service := setup()

	dto := &createDto{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}
	prop := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

	service.On("Create", prop).Return(nil)

	body, err := json.Marshal(dto)
	assert.NoError(t, err)

	// Perform action.
	w := perform("POST", "/api/set", body, router)

	// Test result.
	assert.Equal(t, "/api/set/test.name.1", w.Header().Get("Location"))
	assert.Equal(t, 201, w.Code)
}

func TestCreateConflict(t *testing.T) {
	router, service := setup()

	dto := &createDto{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}
	prop := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

	service.On("Create", prop).Return(apperrors.NewConflict(reflect.TypeOf(prop), "name", prop.Name))

	body, err := json.Marshal(dto)
	assert.NoError(t, err)

	// Perform action.
	w := perform("POST", "/api/set", body, router)

	// Test result.
	assert.Equal(t, 409, w.Code)
}

func TestCreateSyntacticInvalidRequestJSON(t *testing.T) {
	router, service := setup()

	prop := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

	service.On("Create", prop).Return(nil)

	body := []byte(`{"name": "TestCreateDto" "description": "Test Description Dto" "value": "122345"}`)

	// Perform action.
	w := perform("POST", "/api/set", body, router)

	// Test result.
	assert.Equal(t, 400, w.Code)
}

func TestCreateUnexpected(t *testing.T) {
	router, service := setup()

	dto := &createDto{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

	prop := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

	service.On("Create", prop).Return(errors.New("unexpected"))

	body, err := json.Marshal(dto)
	assert.NoError(t, err)

	// Perform action.
	w := perform("POST", "/api/set", body, router)

	// Test result.
	assert.Equal(t, 500, w.Code)
}

func TestReadAll(t *testing.T) {
	router, service := setup()

	// Mock service return.
	properties := make([]*model.PropertySet, 3)
	for i := 0; i < 3; i++ {
		properties[i] = &model.PropertySet{Name: "test.name." + strconv.Itoa(i), Values: []string{"test.value." + strconv.Itoa(i) + ".1", "test.value." + strconv.Itoa(i) + ".2"}}
	}

	service.On("ReadAll").Return(properties, nil)

	// Perform action.
	w := perform("GET", "/api/set", nil, router)

	// Test result.
	assert.Equal(t, 200, w.Code)
}

func TestReadAllUnexpected(t *testing.T) {
	router, service := setup()

	// Mock service return.
	properties := make([]*model.PropertySet, 3)
	for i := 0; i < 3; i++ {
		properties[i] = &model.PropertySet{Name: "test.name." + strconv.Itoa(i), Values: []string{"test.value." + strconv.Itoa(i) + ".1", "test.value." + strconv.Itoa(i) + ".2"}}
	}

	service.On("ReadAll").Return(properties, errors.New("unexpected"))

	// Perform action.
	w := perform("GET", "/api/set", nil, router)

	// Test result.
	assert.Equal(t, 500, w.Code)
}

func TestRead(t *testing.T) {
	router, service := setup()

	// Mock service return.
	property := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

	service.On("FindByID", "test.name.1").Return(property, nil)

	// Perform action.
	w := perform("GET", "/api/set/test.name.1", nil, router)

	// Test result
	assert.Equal(t, 200, w.Code)
}

func TestReadNotFound(t *testing.T) {
	router, service := setup()

	// Mock service return.
	service.On("FindByID", "test.name.1").Return(&model.PropertySet{}, apperrors.NewEntityNotFound(reflect.TypeOf(&model.PropertySet{}), "test.name.1"))

	// Perform action.
	w := perform("GET", "/api/set/test.name.1", nil, router)

	// Test result.
	assert.Equal(t, 404, w.Code)
}

func TestReadUnexpected(t *testing.T) {
	router, service := setup()

	// Mock service error.
	service.On("FindByID", "test.name.1").Return(&model.Property{}, errors.New("unexpected"))

	// Perform action.
	w := perform("GET", "/api/set/test.name.1", nil, router)

	// Test result.
	assert.Equal(t, 500, w.Code)
}

func TestUpdate(t *testing.T) {
	router, service := setup()

	dto := &updateDto{Values: []string{"test.value.1.1", "test.value.1.2"}}

	prop := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

	service.On("Update", prop).Return(nil)

	body, err := json.Marshal(dto)
	assert.NoError(t, err)

	// Perform action.
	w := perform("PUT", "/api/set/test.name.1", body, router)

	// Test result.
	assert.Equal(t, 200, w.Code)
}

func TestUpdateSyntacticInvalidRequestJSON(t *testing.T) {
	router, service := setup()

	prop := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

	service.On("Update", prop).Return(nil)

	body := []byte(`{"id: "testid", ""name": "TestCreateDto" "description": "Test Description Dto" "value": "122345"}`)

	// Perform action.
	w := perform("PUT", "/api/set/testid", body, router)

	// Test result.
	assert.Equal(t, 400, w.Code)
}

func TestUpdateUnexpected(t *testing.T) {
	router, service := setup()

	dto := &updateDto{Values: []string{"test.value.1.1", "test.value.1.2"}}

	prop := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

	service.On("Update", prop).Return(errors.New("unexpected"))

	body, err := json.Marshal(dto)
	assert.NoError(t, err)

	// Perform action.
	w := perform("PUT", "/api/set/test.name.1", body, router)

	// Test result.
	assert.Equal(t, 500, w.Code)
}

func TestDelete(t *testing.T) {
	router, service := setup()

	// Mock service action.
	service.On("Delete", "TestId").Return(nil)

	// Perform action.
	w := perform("DELETE", "/api/set/TestId", nil, router)

	// Test result.
	assert.Equal(t, 204, w.Code)
}

func TestDeleteNotFound(t *testing.T) {
	router, service := setup()

	// Mock service action.
	service.On("Delete", "TestId").Return(apperrors.NewEntityNotFound(reflect.TypeOf(&model.Property{}), "TestId"))

	// Perform action.
	w := perform("DELETE", "/api/set/TestId", nil, router)

	// Test result.
	assert.Equal(t, 404, w.Code)
}

func TestDeleteUnexpected(t *testing.T) {
	router, service := setup()

	// Mock service error.
	service.On("Delete", "TestId").Return(errors.New("unexpected"))

	// Perform action.
	w := perform("DELETE", "/api/set/TestId", nil, router)

	// Test result.
	assert.Equal(t, 500, w.Code)
}

func setup() (r *gin.Engine, serviceMock *PropertyServiceMock) {
	router := gin.Default()
	router.Use(
		jsonAppErrorHandler(),
	)
	api := router.Group("/api")

	service := new(PropertyServiceMock)
	controller := New(service).Controller
	controller.Register(api)

	return router, service
}

func perform(method string, uri string, body []byte, router *gin.Engine) (rr *httptest.ResponseRecorder) {
	return performWithHeaders(method, uri, body, router, nil)
}

func performWithHeaders(method string, uri string, body []byte, router *gin.Engine, headers map[string]string) (rr *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()

	var content bytes.Buffer
	if body != nil {
		content = *bytes.NewBuffer(body)
	}

	req, _ := http.NewRequest(method, uri, &content)
	if headers != nil {
		for key, value := range headers {

			req.Header.Set(key, value)
		}
	}
	router.ServeHTTP(w, req)

	return w
}

type PropertyServiceMock struct {
	mock.Mock
}

func (m *PropertyServiceMock) Create(ctx context.Context, property *model.PropertySet) error {
	args := m.Called(property)

	return args.Error(0)
}

func (m *PropertyServiceMock) ReadAll(ctx context.Context) ([]*model.PropertySet, error) {
	args := m.Called()

	return args.Get(0).([]*model.PropertySet), args.Error(1)
}

func (m *PropertyServiceMock) FindByID(ctx context.Context, id string) (*model.PropertySet, error) {
	args := m.Called(id)

	return args.Get(0).(*model.PropertySet), args.Error(1)
}

func (m *PropertyServiceMock) Delete(ctx context.Context, id string) error {
	args := m.Called(id)

	return args.Error(0)
}

func (m *PropertyServiceMock) Update(ctx context.Context, property *model.PropertySet) error {
	args := m.Called(property)

	return args.Error(0)
}

func jsonAppErrorHandler() gin.HandlerFunc {
	return handle(gin.ErrorTypeAny)
}

func handle(errType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedErrors := c.Errors

		if len(detectedErrors) > 0 {
			err := detectedErrors[0].Err

			switch err.(type) {
			case *apperrors.Error:
				oError := err.(*apperrors.Error)
				c.AbortWithError(oError.Code, oError)
			default:
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			return
		}
	}
}
