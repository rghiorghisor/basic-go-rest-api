package storage

import (
	"github.com/asdine/storm/v3"
	"github.com/rghiorghisor/basic-go-rest-api/config"
	"github.com/rghiorghisor/basic-go-rest-api/logger"
	property_bolt "github.com/rghiorghisor/basic-go-rest-api/property/gateway/storage/bolt"
	propertyset_bolt "github.com/rghiorghisor/basic-go-rest-api/propertyset/gateway/storage/bolt"
	"github.com/rghiorghisor/basic-go-rest-api/util"
)

type boltFactory struct {
}

func newBoltFactory() factory {
	return new(boltFactory)
}

func (f *boltFactory) id() string {
	return "local"
}

func (f *boltFactory) init(storage *Storage, config *config.StorageConfiguration) error {
	boltConfig := config.BoltDbConfiguration
	logger.Main.Info("Setting up local(storm) storage...")

	// Create db connection.
	dbt, err := f.connectStorm(boltConfig)
	if err != nil {
		return err
	}

	// Setup repositories.
	storage.PropertyRepository = property_bolt.New(dbt)
	storage.PropertySetRepository = propertyset_bolt.New(dbt)

	// Add here any new repository...

	return nil
}

func (f *boltFactory) connectStorm(config *config.BoltDbConfiguration) (*storm.DB, error) {
	util.CreateParentFolder(config.Name)

	return storm.Open(config.Name)
}
