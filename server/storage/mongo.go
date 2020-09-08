package storage

import (
	"context"
	"log"
	"time"

	"github.com/rghiorghisor/basic-go-rest-api/config"
	"github.com/rghiorghisor/basic-go-rest-api/logger"
	property_mongo "github.com/rghiorghisor/basic-go-rest-api/property/gateway/storage/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoFactory struct {
}

func newMongoFactory() factory {
	return new(mongoFactory)
}

func (f *mongoFactory) id() string {
	return "mongo"
}

func (f *mongoFactory) init(storage *Storage, config *config.StorageConfiguration) error {
	mongoConfig := config.DbConfiguration
	logger.Main.Info("Setting up mongoDb storage...")

	// Create db connection.
	db := f.connect(mongoConfig)

	// Setup repositories.
	storage.PropertyRepository = property_mongo.New(db, mongoConfig.PropertiesCollectionName)

	// Add here any new repository...

	return nil
}

func (f *mongoFactory) connect(dbConfiguration *config.MongoDbConfiguration) *mongo.Database {
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
