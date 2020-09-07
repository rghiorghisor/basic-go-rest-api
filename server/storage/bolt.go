package storage

import (
	"github.com/asdine/storm/v3"
	"github.com/rghiorghisor/basic-go-rest-api/config"
	"github.com/rghiorghisor/basic-go-rest-api/logger"
	property_bolt "github.com/rghiorghisor/basic-go-rest-api/property/gateway/storage/bbolt"
	"github.com/rghiorghisor/basic-go-rest-api/util"
)

func initBolt(storage *Storage, config *config.BoltDbConfiguration) error {
	logger.Main.Info("Setting up local(storm) storage...")

	// Create db connection.
	dbt, err := connectStorm(config)
	if err != nil {
		return err
	}

	// Setup repositories.
	storage.PropertyRepository = property_bolt.New(dbt)

	// Add here any new repository...

	return nil
}

func connectStorm(config *config.BoltDbConfiguration) (*storm.DB, error) {
	util.CreateParentFolder(config.Name)

	return storm.Open(config.Name)
}
