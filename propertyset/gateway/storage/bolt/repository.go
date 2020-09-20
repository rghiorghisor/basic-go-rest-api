package bolt

import (
	"context"
	"reflect"

	"github.com/asdine/storm/v3"
	"github.com/rghiorghisor/basic-go-rest-api/errors"
	"github.com/rghiorghisor/basic-go-rest-api/model"
	"github.com/rghiorghisor/basic-go-rest-api/propertyset/gateway/storage"
)

// PropertySetRepository is a representation of the property repository for Bolt DBs.
type PropertySetRepository struct {
	db *storm.DB
}

type propertySetDto struct {
	Name   string   `storm:"id,unique"`
	Values []string `bson:"values"`
}

// New retrieves a new repository object ready to be used.
func New(db *storm.DB) storage.Repository {
	repo := &PropertySetRepository{
		db: db,
	}
	db.Init(&propertySetDto{})

	return repo
}

// Create a new entry based on the provided property.
func (repository PropertySetRepository) Create(ctx context.Context, propSet *model.PropertySet) error {
	dto := convertToDto(propSet)

	tx, err := repository.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var found propertySetDto

	e := tx.One("Name", propSet.Name, &found)
	if e == nil {
		return errors.NewConflict(reflect.TypeOf((*model.PropertySet)(nil)).Elem(), "name", propSet.Name)
	}

	if e != storm.ErrNotFound {
		return e
	}

	if err := tx.Save(dto); err != nil {
		return err
	}

	return tx.Commit()
}

// ReadAll retrieves all available properties.
func (repository PropertySetRepository) ReadAll(ctx context.Context) ([]*model.PropertySet, error) {
	var propSets []propertySetDto
	if err := repository.db.All(&propSets); err != nil {
		return nil, err
	}

	return convertDtosToModel(propSets), nil
}

// FindByID retrieves the property matching the given id if such a property
// exists; otherwise will return a not found error.
func (repository PropertySetRepository) FindByID(context context.Context, id string) (*model.PropertySet, error) {
	var dto propertySetDto
	err := repository.db.One("Name", id, &dto)

	if storm.ErrNotFound == err {
		return nil, errors.NewEntityNotFound(model.PropertySet{}, id)
	}

	if err != nil {
		return nil, err
	}

	return convertToModel(&dto), nil
}

// FindByName retrieves the property matching the given name if such a property
// exists; otherwise will return a not found error.
func (repository PropertySetRepository) FindByName(context context.Context, name string) (*model.PropertySet, error) {
	var dto propertySetDto
	err := repository.db.One("Name", name, &dto)

	if storm.ErrNotFound == err {
		return nil, errors.NewEntityNotFound(model.PropertySet{}, name)
	}

	if err != nil {
		return nil, err
	}

	return convertToModel(&dto), nil
}

// Delete the property with the given id.
func (repository PropertySetRepository) Delete(context context.Context, id string) error {
	var dto propertySetDto
	err := repository.db.One("Name", id, &dto)

	if storm.ErrNotFound == err {
		return errors.NewEntityNotFound(model.PropertySet{}, id)
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
func (repository PropertySetRepository) Update(ctx context.Context, property *model.PropertySet) error {
	dto := convertToDto(property)
	err := repository.db.Update(dto)

	if storm.ErrNotFound == err {
		return errors.NewEntityNotFound(model.PropertySet{}, property.Name)
	}

	if err != nil {
		return err
	}

	return nil
}

func convertToDto(property *model.PropertySet) *propertySetDto {
	return &propertySetDto{
		Name:   property.Name,
		Values: property.Values,
	}
}

func convertDtosToModel(dtos []propertySetDto) []*model.PropertySet {
	result := make([]*model.PropertySet, len(dtos))

	for index, dto := range dtos {
		result[index] = convertToModel(&dto)
	}

	return result
}

func convertToModel(dto *propertySetDto) *model.PropertySet {
	return &model.PropertySet{
		Name:   dto.Name,
		Values: dto.Values,
	}
}
