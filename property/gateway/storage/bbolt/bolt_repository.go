package bbolt

import (
	"context"
	"reflect"

	"github.com/asdine/storm/v3"
	"github.com/google/uuid"
	"github.com/rghiorghisor/basic-go-rest-api/errors"
	"github.com/rghiorghisor/basic-go-rest-api/model"
	"github.com/rghiorghisor/basic-go-rest-api/property/gateway/storage"
)

// PropertyRepository is a representation of the property repository for Bolt DBs.
type PropertyRepository struct {
	db *storm.DB
}

type propertyDto struct {
	ID          string `storm:"id"`
	Name        string `storm:"unique"`
	Description string `bson:"description"`
	Value       string `bson:"value"`
}

// New retrieves a new repository object ready to be used.
func New(db *storm.DB) storage.Repository {
	repo := &PropertyRepository{
		db: db,
	}
	db.Init(&propertyDto{})

	return repo
}

// Create a new entry based on the provided property.
func (repository PropertyRepository) Create(ctx context.Context, property *model.Property) error {
	id := uuid.New().String()
	property.ID = id

	dto := convertToDto(property)

	err := repository.db.Save(dto)
	if err != nil {
		return err
	}

	return nil
}

// ReadAll retrieves all available properties.
func (repository PropertyRepository) ReadAll(ctx context.Context) ([]*model.Property, error) {
	var properties []propertyDto
	if err := repository.db.All(&properties); err != nil {
		return nil, err
	}

	return convertDtosToModel(properties), nil
}

// FindByID retrieves the property matching the given id if such a property
// exists; otherwise will return a not found error.
func (repository PropertyRepository) FindByID(context context.Context, id string) (*model.Property, error) {
	var dto propertyDto
	err := repository.db.One("ID", id, &dto)

	if storm.ErrNotFound == err {
		return nil, errors.NewEntityNotFound(reflect.TypeOf((*model.Property)(nil)).Elem(), id)
	}

	if err != nil {
		return nil, err
	}

	return convertToModel(&dto), nil
}

// FindByName retrieves the property matching the given name if such a property
// exists; otherwise will return a not found error.
func (repository PropertyRepository) FindByName(context context.Context, name string) (*model.Property, error) {
	var dto propertyDto
	err := repository.db.One("Name", name, &dto)

	if storm.ErrNotFound == err {
		return nil, errors.NewEntityNotFound(reflect.TypeOf((*model.Property)(nil)).Elem(), name)
	}

	if err != nil {
		return nil, err
	}

	return convertToModel(&dto), nil
}

// Delete the property with the given id.
func (repository PropertyRepository) Delete(context context.Context, id string) error {
	var dto propertyDto
	err := repository.db.One("ID", id, &dto)

	if storm.ErrNotFound == err {
		return errors.NewEntityNotFound(reflect.TypeOf((*model.Property)(nil)).Elem(), id)
	}

	if err != nil {
		return err
	}

	if err := repository.db.DeleteStruct(&dto); err != nil {
		return err
	}

	return nil
}

// Update all fields of the given property.
func (repository PropertyRepository) Update(ctx context.Context, property *model.Property) error {
	dto := convertToDto(property)
	err := repository.db.Update(dto)

	if storm.ErrNotFound == err {
		return errors.NewEntityNotFound(reflect.TypeOf((*model.Property)(nil)).Elem(), property.ID)
	}

	if err != nil {
		return err
	}

	return nil
}

func convertToDto(property *model.Property) *propertyDto {
	return &propertyDto{
		ID:          property.ID,
		Name:        property.Name,
		Description: property.Description,
		Value:       property.Value,
	}
}

func convertDtosToModel(dtos []propertyDto) []*model.Property {
	result := make([]*model.Property, len(dtos))

	for index, dto := range dtos {
		result[index] = convertToModel(&dto)
	}

	return result
}

func convertToModel(dto *propertyDto) *model.Property {
	return &model.Property{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Value:       dto.Value,
	}
}
