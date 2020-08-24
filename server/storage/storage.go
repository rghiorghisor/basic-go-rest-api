package storage

import (
	"context"
	"log"
	"time"

	"github.com/rghiorghisor/basic-go-rest-api/config"
	pstorage "github.com/rghiorghisor/basic-go-rest-api/property/gateway/storage"
	pgmongo "github.com/rghiorghisor/basic-go-rest-api/property/gateway/storage/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Storage structure contains all repositories.
type Storage struct {
	PropertyRepository pstorage.Repository
}

// NewStorage returns a bare-bone storage.
func NewStorage() *Storage {
	return &Storage{}
}

// SetupStorage prepares the MongoDB connection based on the provided
// configuration and retrieves a handle for the database.
//
// Besides connecting to the database it also prepares repositories based on any
// collection names.
func (storage *Storage) SetupStorage(dbConfiguration *config.MongoDbConfiguration) {
	// Create db connection.
	db := connect(dbConfiguration)

	// Create property repository.
	propCollectionName := dbConfiguration.PropertiesCollectionName
	storage.PropertyRepository = pgmongo.NewMongoPropertyRepository(db, propCollectionName)
}

func connect(dbConfiguration *config.MongoDbConfiguration) *mongo.Database {
	uri := dbConfiguration.URI
	dbName := dbConfiguration.Name

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Cannot establish mongoDB connection:")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(dbName)
}
