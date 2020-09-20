package service

import (
	"context"

	"github.com/rghiorghisor/basic-go-rest-api/errors"
	"github.com/rghiorghisor/basic-go-rest-api/model"
	"github.com/rghiorghisor/basic-go-rest-api/propertyset"
	"github.com/rghiorghisor/basic-go-rest-api/propertyset/gateway/storage"
	serverstorage "github.com/rghiorghisor/basic-go-rest-api/server/storage"
)

// PropertySetService defines the service handling property sets operations.
type PropertySetService struct {
	repository storage.Repository
}

// New creates a PropertySetService.
//
// As this service needs access to a repository to perform action, it is the
// responsibility of the service to get the correct repo from the storage parameter.
func New(storage *serverstorage.Storage) propertyset.Service {
	return PropertySetService{
		repository: storage.PropertySetRepository,
	}
}

// Create processes a new property set and adds it to the repository.
func (service PropertySetService) Create(ctx context.Context, prop *model.PropertySet) error {
	return service.repository.Create(ctx, prop)
}

// ReadAll retrieves all available property sets.
func (service PropertySetService) ReadAll(ctx context.Context) ([]*model.PropertySet, error) {
	return service.repository.ReadAll(ctx)
}

// FindByID retrieves the property set matching the given id if such a property set
// exists; otherwise will return a not found error.
func (service PropertySetService) FindByID(ctx context.Context, id string) (*model.PropertySet, error) {
	foundProp, err := service.repository.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if foundProp == nil {
		return nil, errors.NewEntityNotFound(foundProp, id)
	}

	return foundProp, nil
}

// FindValuesByID retrieves all values associated with the set identified by the
// given parameter.
func (service PropertySetService) FindValuesByID(ctx context.Context, id string) ([]string, error) {
	foundSet, err := service.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return foundSet.Values, nil
}

// Delete the property set with the given id.
func (service PropertySetService) Delete(ctx context.Context, id string) error {
	return service.repository.Delete(ctx, id)
}

// Update all fields of the given property set.
func (service PropertySetService) Update(ctx context.Context, prop *model.PropertySet) error {
	return service.repository.Update(ctx, prop)
}
