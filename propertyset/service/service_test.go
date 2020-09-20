package service

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"testing"

	apperrors "github.com/rghiorghisor/basic-go-rest-api/errors"
	"github.com/rghiorghisor/basic-go-rest-api/model"
	"github.com/rghiorghisor/basic-go-rest-api/propertyset"
	"github.com/rghiorghisor/basic-go-rest-api/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	srv, repo := setup()

	toCreate := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

	repo.On("Create", toCreate).Return(nil)

	ctx := context.Background()
	actualErr := srv.Create(ctx, toCreate)

	assert.Nil(t, actualErr)
}

func TestCreateConflict(t *testing.T) {
	srv, repo := setup()

	toCreate := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

	repo.On("Create", toCreate).Return(apperrors.NewConflict(reflect.TypeOf(toCreate), "name", "test.name.1"))

	ctx := context.Background()
	actualErr := srv.Create(ctx, toCreate)

	assert.Equal(t, apperrors.NewConflict(reflect.TypeOf(toCreate), "name", "test.name.1"), actualErr)
}

func TestReadAll(t *testing.T) {
	srv, repo := setup()

	properties := make([]*model.PropertySet, 3)
	for i := 0; i < 3; i++ {
		properties[i] = &model.PropertySet{Name: "test.name." + strconv.Itoa(i), Values: []string{"test.value." + strconv.Itoa(i) + ".1", "test.value." + strconv.Itoa(i) + ".2"}}
	}

	repo.On("ReadAll").Return(properties, nil)

	ctx := context.Background()
	actual, err := srv.ReadAll(ctx)

	assert.Nil(t, err)
	assert.Equal(t, properties, actual)
}

func TestFindByID(t *testing.T) {
	srv, repo := setup()

	found := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}
	repo.On("FindByID", found.Name).Return(found, nil)

	ctx := context.Background()
	actual, err := srv.FindByID(ctx, found.Name)

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
	assert.Equal(t, apperrors.NewEntityNotFound(actual, "testid"), err)
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

	toUpdate := &model.PropertySet{Name: "test.name.1", Values: []string{"test.value.1.1", "test.value.1.2"}}

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

func setup() (service propertyset.Service, repo *PropertyRepositoryMock) {
	repoMock := new(PropertyRepositoryMock)
	storage := &storage.Storage{PropertySetRepository: repoMock}
	service = New(storage)

	return service, repoMock
}

type PropertyRepositoryMock struct {
	mock.Mock
}

func (m PropertyRepositoryMock) Create(ctx context.Context, property *model.PropertySet) error {
	args := m.Called(property)

	return args.Error(0)
}

func (m PropertyRepositoryMock) ReadAll(ctx context.Context) ([]*model.PropertySet, error) {
	args := m.Called()

	return args.Get(0).([]*model.PropertySet), args.Error(1)
}

func (m PropertyRepositoryMock) FindByID(context context.Context, id string) (*model.PropertySet, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.PropertySet), args.Error(1)
}

func (m PropertyRepositoryMock) FindByName(context context.Context, name string) (*model.PropertySet, error) {
	args := m.Called(name)

	var q *model.PropertySet
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	q = args.Get(0).(*model.PropertySet)

	return q, args.Error(1)
}

func (m PropertyRepositoryMock) Delete(context context.Context, id string) error {
	args := m.Called(id)

	return args.Error(0)

}

func (m PropertyRepositoryMock) Update(ctx context.Context, property *model.PropertySet) error {
	args := m.Called(property)

	return args.Error(0)
}
