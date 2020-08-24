package service

import (
	"context"
	"reflect"

	"github.com/rghiorghisor/basic-go-rest-api/errors"
	"github.com/rghiorghisor/basic-go-rest-api/model"
	"github.com/rghiorghisor/basic-go-rest-api/property/gateway/storage"
	serverstorage "github.com/rghiorghisor/basic-go-rest-api/server/storage"
)

// PropertyService defines the service handling property operations.
type PropertyService struct {
	repository storage.Repository ``
}

// NewPropertyService creates a PropertyService.
//
// As this service needs access to a repository to perform action, it is the
// responsibility of the service to get the correct repo from the storage parameter.
func NewPropertyService(storage *serverstorage.Storage) *PropertyService {
	return &PropertyService{
		repository: storage.PropertyRepository,
	}
}

// Create processes a new property and adds it to the repository.
func (service PropertyService) Create(ctx context.Context, prop *model.Property) error {
	foundProp, _ := service.repository.FindByName(ctx, prop.Name)
	if foundProp != nil {
		return errors.NewConflict(reflect.TypeOf(foundProp), "name", prop.Name)
	}

	return service.repository.Create(ctx, prop)
}

// ReadAll retrieves all available properties.
func (service PropertyService) ReadAll(ctx context.Context) ([]*model.Property, error) {
	return service.repository.ReadAll(ctx)
}

// FindByID retrieves the property matching the given id if such a property
// exists; otherwise will return a not found error.
func (service PropertyService) FindByID(ctx context.Context, id string) (*model.Property, error) {
	foundProp, err := service.repository.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if foundProp == nil {
		return nil, errors.NewEntityNotFound(reflect.TypeOf(foundProp), id)
	}

	return foundProp, nil
}

// Delete the property with the given id.
func (service PropertyService) Delete(ctx context.Context, id string) error {
	return service.repository.Delete(ctx, id)
}

// Update all fields of the given property.
func (service PropertyService) Update(ctx context.Context, property *model.Property) error {
	return service.repository.Update(ctx, property)
}
