package mongo

import (
	"context"
	"reflect"
	"strings"

	"github.com/rghiorghisor/basic-go-rest-api/errors"
	"github.com/rghiorghisor/basic-go-rest-api/model"
	"github.com/rghiorghisor/basic-go-rest-api/propertyset/gateway/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type propertySetDto struct {
	Name   string   `bson:"_id"`
	Values []string `bson:"values"`
}

// PropertySetRepository is a representation of the property repository for
// a mongo DBs.
type PropertySetRepository struct {
	dbCollection *mongo.Collection
}

// New retrieves a new repository object ready to be used.
func New(db *mongo.Database, collectionName string) storage.Repository {
	return &PropertySetRepository{
		dbCollection: db.Collection("set_collection"),
	}
}

// Create a new entry based on the provided property.
func (repository PropertySetRepository) Create(ctx context.Context, property *model.PropertySet) error {
	dto := convertToDto(property)

	_, err := repository.dbCollection.InsertOne(ctx, dto)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			return errors.NewConflict(reflect.TypeOf(property), "name", property.Name)
		}

		return err
	}

	return nil
}

// ReadAll retrieves all available properties.
func (repository PropertySetRepository) ReadAll(ctx context.Context) ([]*model.PropertySet, error) {
	cursor, error := repository.dbCollection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)

	if error != nil {
		return nil, error
	}

	result := make([]*propertySetDto, 0)

	for cursor.Next(ctx) {
		dto := new(propertySetDto)
		err := cursor.Decode(dto)
		if err != nil {
			return nil, err
		}

		result = append(result, dto)
	}

	return convertDtosToModel(result), nil
}

// FindByID retrieves the property matching the given id if such a property
// exists; otherwise will return a not found error.
func (repository PropertySetRepository) FindByID(context context.Context, id string) (*model.PropertySet, error) {
	objID := id

	result := new(propertySetDto)
	err := repository.dbCollection.FindOne(context, bson.M{
		"_id": objID,
	}).Decode(result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewEntityNotFound(reflect.TypeOf((*model.Property)(nil)).Elem(), id)
		}

		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return convertToModel(result), nil
}

func (repository PropertySetRepository) findAllBy(ctx context.Context, queryValues *map[string]string) ([]*model.PropertySet, error) {
	filter := bson.M{}

	for k, v := range *queryValues {
		filter[k] = v
	}

	cursor, err := repository.dbCollection.Find(ctx, filter)
	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]*propertySetDto, 0)

	for cursor.Next(ctx) {
		dto := new(propertySetDto)
		err := cursor.Decode(dto)
		if err != nil {
			return nil, err
		}

		result = append(result, dto)
	}

	return convertDtosToModel(result), nil
}

func (repository PropertySetRepository) findOneBy(context context.Context, queryValues *map[string]string) (*model.PropertySet, error) {
	result := new(propertySetDto)
	filter := bson.M{}

	for k, v := range *queryValues {
		filter[k] = v
	}

	error := repository.dbCollection.FindOne(context, filter).Decode(result)

	if error != nil {
		return nil, error
	}

	if result == nil {
		return nil, nil
	}

	return convertToModel(result), nil
}

// Delete the property with the given id.
func (repository PropertySetRepository) Delete(context context.Context, id string) error {
	objID := id

	_, err := repository.dbCollection.DeleteOne(context, bson.M{
		"_id": objID})

	return err
}

// Update all fields of the given property.
func (repository PropertySetRepository) Update(ctx context.Context, property *model.PropertySet) error {
	_, err := repository.dbCollection.UpdateOne(ctx,
		bson.M{"_id": property.Name},
		bson.D{primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "values", Value: property.Values},
		}}})

	return err
}

func convertToDto(property *model.PropertySet) *propertySetDto {
	return &propertySetDto{
		Name:   property.Name,
		Values: property.Values,
	}
}

func convertDtosToModel(dtos []*propertySetDto) []*model.PropertySet {
	result := make([]*model.PropertySet, len(dtos))

	for index, dto := range dtos {
		result[index] = convertToModel(dto)
	}

	return result
}

func convertToModel(dto *propertySetDto) *model.PropertySet {
	return &model.PropertySet{
		Name:   dto.Name,
		Values: dto.Values,
	}
}
