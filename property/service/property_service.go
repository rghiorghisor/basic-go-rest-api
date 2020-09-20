package service

import (
	"context"
	"reflect"

	"github.com/rghiorghisor/basic-go-rest-api/errors"
	"github.com/rghiorghisor/basic-go-rest-api/model"
	"github.com/rghiorghisor/basic-go-rest-api/property"
	"github.com/rghiorghisor/basic-go-rest-api/property/gateway/storage"
	"github.com/rghiorghisor/basic-go-rest-api/propertyset"
	serverstorage "github.com/rghiorghisor/basic-go-rest-api/server/storage"
)

// PropertyService defines the service handling property operations.
type PropertyService struct {
	validators validators
	repository storage.Repository

	setService propertyset.Service
}

// New creates a PropertyService.
//
// As this service needs access to a repository to perform action, it is the
// responsibility of the service to get the correct repo from the storage parameter.
func New(storage *serverstorage.Storage, setService propertyset.Service) property.Service {
	return PropertyService{
		validators: newValidators(),
		repository: storage.PropertyRepository,
		setService: setService,
	}
}

// Create processes a new property and adds it to the repository.
func (service PropertyService) Create(ctx context.Context, prop *model.Property) error {
	if err := service.validators.check(prop); err != nil {
		return err
	}

	foundProp, _ := service.repository.FindByName(ctx, prop.Name)
	if foundProp != nil {
		return errors.NewConflict(reflect.TypeOf(foundProp), "name", prop.Name)
	}

	return service.repository.Create(ctx, prop)
}

// ReadAll retrieves all available properties, based on the provided query.
//
// The Query.Set defines the set of properties to be used. In case such a set
// is defined, the names from the set will be used to filter the results;
// otherwise, all properties are retrieved.
func (service PropertyService) ReadAll(ctx context.Context, query property.Query) ([]*model.Property, error) {
	var filterValues []string
	var err error
	if query.HasSet() {
		filterValues, err = service.setService.FindValuesByID(ctx, query.GetSet())
		if err != nil {
			return nil, err
		}

		return service.repository.ReadAllFiltered(ctx, filterValues)
	}

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
		return nil, errors.NewEntityNotFound(model.Property{}, id)
	}

	return foundProp, nil
}

// Delete the property with the given id.
func (service PropertyService) Delete(ctx context.Context, id string) error {
	return service.repository.Delete(ctx, id)
}

// Update all fields of the given property.
func (service PropertyService) Update(ctx context.Context, prop *model.Property) error {
	if err := service.validators.check(prop); err != nil {
		return err
	}

	return service.repository.Update(ctx, prop)
}
