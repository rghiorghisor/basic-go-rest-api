package storage

import (
	"strings"

	"github.com/rghiorghisor/basic-go-rest-api/config"
	"github.com/rghiorghisor/basic-go-rest-api/logger"
	property "github.com/rghiorghisor/basic-go-rest-api/property/gateway/storage"
	propertyset "github.com/rghiorghisor/basic-go-rest-api/propertyset/gateway/storage"
)

// Storage structure contains all repositories.
type Storage struct {
	factories             []func() factory
	defaultFactory        func() factory
	PropertyRepository    property.Repository
	PropertySetRepository propertyset.Repository
}

type factory interface {
	id() string
	init(storage *Storage, config *config.StorageConfiguration) error
}

// New returns a bare-bone storage.
func New() *Storage {
	return &Storage{
		factories:      []func() factory{newBoltFactory, newMongoFactory},
		defaultFactory: newBoltFactory,
	}
}

// SetupStorage prepares the repository connections based on the provided
// configuration and retrieves a handle for the database.
//
// Besides connecting to the database it also prepares repositories based on any
// collection names.
func (storage *Storage) SetupStorage(config *config.StorageConfiguration) error {
	var factory factory
	for _, ff := range storage.factories {
		f := ff()
		if !checkConfig(f, config.Type) {
			continue
		}

		factory = f
		break
	}

	if factory == nil {
		factory = storage.defaultFactory()
		logger.Main.Infof("Unknown storage type '%s'. Using default '%s'.\n", config.Type, factory.id())
	}

	err := factory.init(storage, config)

	return err
}

func checkConfig(f factory, storageType string) bool {
	return strings.EqualFold(f.id(), storageType)
}
