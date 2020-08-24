package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	apperrors "github.com/rghiorghisor/basic-go-rest-api/errors"
	"github.com/rghiorghisor/basic-go-rest-api/model"
	"github.com/rghiorghisor/basic-go-rest-api/property"
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
	actual, err := srv.ReadAll(ctx)

	assert.Nil(t, err)
	assert.Equal(t, properties, actual)
}

func TestFindById(t *testing.T) {
	srv, repo := setup()

	found := &model.Property{
		ID:    "TestId",
		Name:  "TestName",
		Value: "TestValue"}

	repo.On("FindById", found.ID).Return(found, nil)

	ctx := context.Background()
	actual, err := srv.FindById(ctx, found.ID)

	assert.Nil(t, err)
	assert.Equal(t, found, actual)
}

func TestFindByIdNotFound(t *testing.T) {
	srv, repo := setup()

	notFoundID := "testid"
	repo.On("FindById", notFoundID).Return(nil, nil)

	ctx := context.Background()
	actual, err := srv.FindById(ctx, notFoundID)

	assert.Nil(t, actual)
	assert.Equal(t, apperrors.NewEntityNotFound(reflect.TypeOf(actual), "testid"), err)
}

func TestFindByIdUnexpected(t *testing.T) {
	srv, repo := setup()

	notFoundID := "TestId"
	expectedError := errors.New("unexpected")
	repo.On("FindById", notFoundID).Return(nil, expectedError)

	ctx := context.Background()
	actual, err := srv.FindById(ctx, notFoundID)

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
	storage := &storage.Storage{repoMock}
	service = NewPropertyService(storage)

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

func (m PropertyRepositoryMock) FindById(context context.Context, id string) (*model.Property, error) {
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
