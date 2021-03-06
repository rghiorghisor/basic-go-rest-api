package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	apperrors "github.com/rghiorghisor/basic-go-rest-api/errors"
	"github.com/rghiorghisor/basic-go-rest-api/model"
	"github.com/rghiorghisor/basic-go-rest-api/property"
	set_service "github.com/rghiorghisor/basic-go-rest-api/propertyset/service"
	"github.com/rghiorghisor/basic-go-rest-api/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	srv, repo := setup()

	toCreate := &model.Property{
		Name:  "TestName",
		Value: "TestValue"}

	repo.On("FindByName", toCreate.Name).Return(nil, nil)
	repo.On("Create", toCreate).Return(nil)

	ctx := context.Background()
	actualErr := srv.Create(ctx, toCreate)

	assert.Nil(t, actualErr)
}

func TestCreateConflict(t *testing.T) {
	srv, repo := setup()

	toCreate := &model.Property{
		Name:  "TestName",
		Value: "TestValue"}

	found := &model.Property{
		Name:  "TestName",
		Value: "Exiting Value"}

	repo.On("FindByName", toCreate.Name).Return(found, nil)

	ctx := context.Background()
	actualErr := srv.Create(ctx, toCreate)

	assert.Equal(t, apperrors.NewConflict(reflect.TypeOf(found), "name", "TestName"), actualErr)
}

func TestCreateInvalidName(t *testing.T) {
	srv, _ := setup()

	toCreate := &model.Property{
		Name:  "",
		Value: "TestValue"}

	ctx := context.Background()
	actualErr := srv.Create(ctx, toCreate)

	assert.Equal(t, apperrors.NewInvalidEntityEmpty(reflect.TypeOf(model.Property{}), "name"), actualErr)
}

func TestCreateInvalidNameMissing(t *testing.T) {
	srv, _ := setup()

	toCreate := &model.Property{
		Value: "TestValue"}

	ctx := context.Background()
	actualErr := srv.Create(ctx, toCreate)

	assert.Equal(t, apperrors.NewInvalidEntityEmpty(reflect.TypeOf(model.Property{}), "name"), actualErr)
}

func TestCreateInvalidNameSpace(t *testing.T) {
	srv, _ := setup()

	toCreate := &model.Property{
		Name:  "Test Name",
		Value: "TestValue"}

	ctx := context.Background()
	actualErr := srv.Create(ctx, toCreate)

	assert.Equal(t, apperrors.NewInvalidEntityCustom(reflect.TypeOf(model.Property{}), "'name' cannot contain spaces."), actualErr)
}

func TestReadAll(t *testing.T) {
	srv, repo := setup()

	properties := make([]*model.Property, 3)
	for i := 0; i < 3; i++ {
		properties[i] = &model.Property{
			ID:          "Id",
			Name:        "Name test",
			Description: "Description test",
			Value:       "Value test"}
	}

	repo.On("ReadAll").Return(properties, nil)

	ctx := context.Background()
	actual, err := srv.ReadAll(ctx, property.Query{})

	assert.Nil(t, err)
	assert.Equal(t, properties, actual)
}

func TestFindByID(t *testing.T) {
	srv, repo := setup()

	found := &model.Property{
		ID:    "TestId",
		Name:  "TestName",
		Value: "TestValue"}

	repo.On("FindByID", found.ID).Return(found, nil)

	ctx := context.Background()
	actual, err := srv.FindByID(ctx, found.ID)

	assert.Nil(t, err)
	assert.Equal(t, found, actual)
}

func TestFindByIDNotFound(t *testing.T) {
	srv, repo := setup()

	notFoundID := "testid"
	repo.On("FindByID", notFoundID).Return(nil, nil)

	ctx := context.Background()
	actual, err := srv.FindByID(ctx, notFoundID)

	assert.Nil(t, actual)
	assert.Equal(t, apperrors.NewEntityNotFound(model.Property{}, "testid"), err)
}

func TestFindByIDUnexpected(t *testing.T) {
	srv, repo := setup()

	notFoundID := "TestId"
	expectedError := errors.New("unexpected")
	repo.On("FindByID", notFoundID).Return(nil, expectedError)

	ctx := context.Background()
	actual, err := srv.FindByID(ctx, notFoundID)

	assert.Nil(t, actual)
	assert.Equal(t, expectedError, err)
}

func TestUpdate(t *testing.T) {
	srv, repo := setup()

	toUpdate := &model.Property{
		Name:  "TestName",
		Value: "TestValue"}

	repo.On("Update", toUpdate).Return(nil)

	ctx := context.Background()
	err := srv.Update(ctx, toUpdate)

	assert.Nil(t, err)
}

func TestUpdateInvalidName(t *testing.T) {
	srv, _ := setup()

	toUpdate := &model.Property{
		Name:  "",
		Value: "TestValue"}

	ctx := context.Background()
	err := srv.Update(ctx, toUpdate)

	assert.Equal(t, apperrors.NewInvalidEntityEmpty(reflect.TypeOf(model.Property{}), "name"), err)
}

func TestUpdateInvalidNameMissing(t *testing.T) {
	srv, _ := setup()

	toUpdate := &model.Property{
		Value: "TestValue"}

	ctx := context.Background()
	err := srv.Update(ctx, toUpdate)

	assert.Equal(t, apperrors.NewInvalidEntityEmpty(reflect.TypeOf(model.Property{}), "name"), err)
}

func TestUpdateInvalidNameSpace(t *testing.T) {
	srv, _ := setup()

	toUpdate := &model.Property{
		Name:  "Test Name",
		Value: "TestValue"}

	ctx := context.Background()
	err := srv.Update(ctx, toUpdate)

	assert.Equal(t, apperrors.NewInvalidEntityCustom(reflect.TypeOf(model.Property{}), "'name' cannot contain spaces."), err)
}

func TestDelete(t *testing.T) {
	srv, repo := setup()

	toDeleteID := "TestID"

	repo.On("Delete", toDeleteID).Return(nil)

	ctx := context.Background()
	err := srv.Delete(ctx, toDeleteID)

	assert.Nil(t, err)
}

func setup() (service property.Service, repo *PropertyRepositoryMock) {
	repoMock := new(PropertyRepositoryMock)
	storage := &storage.Storage{PropertyRepository: repoMock}

	setService := new(set_service.PropertySetServiceMock)

	service = New(storage, setService)

	return service, repoMock
}

type PropertyRepositoryMock struct {
	mock.Mock
}

func (m PropertyRepositoryMock) Create(ctx context.Context, property *model.Property) error {
	args := m.Called(property)

	return args.Error(0)
}

func (m PropertyRepositoryMock) ReadAll(ctx context.Context) ([]*model.Property, error) {
	args := m.Called()

	return args.Get(0).([]*model.Property), args.Error(1)
}

func (m PropertyRepositoryMock) ReadAllFiltered(ctx context.Context, names []string) ([]*model.Property, error) {
	args := m.Called(names)

	return args.Get(0).([]*model.Property), args.Error(1)
}

func (m PropertyRepositoryMock) FindByID(context context.Context, id string) (*model.Property, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.Property), args.Error(1)
}

func (m PropertyRepositoryMock) FindByName(context context.Context, name string) (*model.Property, error) {
	args := m.Called(name)

	var q *model.Property
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	q = args.Get(0).(*model.Property)

	return q, args.Error(1)
}

func (m PropertyRepositoryMock) Delete(context context.Context, id string) error {
	args := m.Called(id)

	return args.Error(0)

}

func (m PropertyRepositoryMock) Update(ctx context.Context, property *model.Property) error {
	args := m.Called(property)

	return args.Error(0)
}
