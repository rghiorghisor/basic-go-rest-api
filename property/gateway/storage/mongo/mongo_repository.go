package mongo

import (
	"context"
	"reflect"

	"github.com/rghiorghisor/basic-go-rest-api/errors"
	"github.com/rghiorghisor/basic-go-rest-api/model"
	"github.com/rghiorghisor/basic-go-rest-api/property/gateway/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type propertyDto struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Value       string             `bson:"value"`
}

// MongoPropertyRepository is a representation of the property repository for
// a mongo DBs.
type MongoPropertyRepository struct {
	dbCollection *mongo.Collection
}

// NewMongoPropertyRepository retrieves a new repository object ready to be used.
func NewMongoPropertyRepository(db *mongo.Database, collectionName string) storage.Repository {
	return &MongoPropertyRepository{
		dbCollection: db.Collection(collectionName),
	}
}

// Create a new entry based on the provided property.
func (repository MongoPropertyRepository) Create(ctx context.Context, property *model.Property) error {
	dto := convertToDto(property)

	foundProp, _ := repository.FindByName(ctx, property.Name)
	if foundProp != nil {
		return errors.NewConflict(reflect.TypeOf(foundProp), "name", property.Name)
	}

	result, err := repository.dbCollection.InsertOne(ctx, dto)
	if err != nil {
		return err
	}

	property.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return nil
}

// ReadAll retrieves all available properties.
func (repository MongoPropertyRepository) ReadAll(ctx context.Context) ([]*model.Property, error) {
	cursor, error := repository.dbCollection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)

	if error != nil {
		return nil, error
	}

	result := make([]*propertyDto, 0)

	for cursor.Next(ctx) {
		dto := new(propertyDto)
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
func (repository MongoPropertyRepository) FindByID(context context.Context, id string) (*model.Property, error) {
	objID, _ := primitive.ObjectIDFromHex(id)

	result := new(propertyDto)
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

// FindByName retrieves the property matching the given name if such a property
// exists; otherwise will return a not found error.
func (repository MongoPropertyRepository) FindByName(context context.Context, name string) (*model.Property, error) {
	query := make(map[string]string)
	query["name"] = name

	return repository.findOneBy(context, &query)

}

func (repository MongoPropertyRepository) findAllBy(ctx context.Context, queryValues *map[string]string) ([]*model.Property, error) {
	filter := bson.M{}

	for k, v := range *queryValues {
		filter[k] = v
	}

	cursor, err := repository.dbCollection.Find(ctx, filter)
	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]*propertyDto, 0)

	for cursor.Next(ctx) {
		dto := new(propertyDto)
		err := cursor.Decode(dto)
		if err != nil {
			return nil, err
		}

		result = append(result, dto)
	}

	return convertDtosToModel(result), nil
}

func (repository MongoPropertyRepository) findOneBy(context context.Context, queryValues *map[string]string) (*model.Property, error) {
	result := new(propertyDto)
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
func (repository MongoPropertyRepository) Delete(context context.Context, id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)

	_, err := repository.dbCollection.DeleteOne(context, bson.M{
		"_id": objID})

	return err
}

// Update all fields of the given property.
func (repository MongoPropertyRepository) Update(ctx context.Context, property *model.Property) error {
	objID, _ := primitive.ObjectIDFromHex(property.ID)

	_, err := repository.dbCollection.UpdateOne(ctx,
		bson.M{"_id": objID},
		bson.D{primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "name", Value: property.Name},
			primitive.E{Key: "description", Value: property.Description},
			primitive.E{Key: "value", Value: property.Value},
		}}})

	return err
}

func convertToDto(property *model.Property) *propertyDto {
	return &propertyDto{
		Name:        property.Name,
		Description: property.Description,
		Value:       property.Value,
	}
}

func convertDtosToModel(dtos []*propertyDto) []*model.Property {
	result := make([]*model.Property, len(dtos))

	for index, dto := range dtos {
		result[index] = convertToModel(dto)
	}

	return result
}

func convertToModel(dto *propertyDto) *model.Property {
	return &model.Property{
		ID:          dto.ID.Hex(),
		Name:        dto.Name,
		Description: dto.Description,
		Value:       dto.Value,
	}
}
