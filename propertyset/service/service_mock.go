package service

import (
	"context"

	"github.com/rghiorghisor/basic-go-rest-api/model"
	"github.com/stretchr/testify/mock"
)

// PropertySetServiceMock retrieves a new mock for PropertySetService.
type PropertySetServiceMock struct {
	mock.Mock
}

// Create mock function.
func (m *PropertySetServiceMock) Create(ctx context.Context, property *model.PropertySet) error {
	args := m.Called(property)

	return args.Error(0)
}

// ReadAll mock function.
func (m *PropertySetServiceMock) ReadAll(ctx context.Context) ([]*model.PropertySet, error) {
	args := m.Called()

	return args.Get(0).([]*model.PropertySet), args.Error(1)
}

// FindByID mock function.
func (m *PropertySetServiceMock) FindByID(ctx context.Context, id string) (*model.PropertySet, error) {
	args := m.Called(id)

	return args.Get(0).(*model.PropertySet), args.Error(1)
}

// FindValuesByID mock function.
func (m *PropertySetServiceMock) FindValuesByID(ctx context.Context, id string) ([]string, error) {
	args := m.Called(id)

	return args.Get(0).([]string), args.Error(1)
}

// Delete mock function.
func (m *PropertySetServiceMock) Delete(ctx context.Context, id string) error {
	args := m.Called(id)

	return args.Error(0)
}

// Update mock function.
func (m *PropertySetServiceMock) Update(ctx context.Context, property *model.PropertySet) error {
	args := m.Called(property)

	return args.Error(0)
}
