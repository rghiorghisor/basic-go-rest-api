package storage

import (
	"strings"

	"github.com/rghiorghisor/basic-go-rest-api/config"
	"github.com/rghiorghisor/basic-go-rest-api/logger"
	pstorage "github.com/rghiorghisor/basic-go-rest-api/property/gateway/storage"
)

// Storage structure contains all repositories.
type Storage struct {
	PropertyRepository pstorage.Repository
}

// NewStorage returns a bare-bone storage.
func NewStorage() *Storage {
	return &Storage{}
}

// SetupStorage prepares the repository connections based on the provided
// configuration and retrieves a handle for the database.
//
// Besides connecting to the database it also prepares repositories based on any
// collection names.
func (storage *Storage) SetupStorage(config *config.StorageConfiguration) {
	if strings.EqualFold("mongo", config.Type) {
		// Check the initMongo function to add any new mongoDB repository.
		initMongo(storage, config.DbConfiguration)

		return
	}

	if !strings.EqualFold("local", config.Type) {
		logger.Main.Infof("Unknown storage type '%s'. Using default '%s'.\r\n", config.Type, "local")
	}

	// Check the initBolt function to add any new local repository.
	initBolt(storage, config.BoltDbConfiguration)

}
