package config

import (
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestLoad(t *testing.T) {
	appConfiguration := setupAndLoad("test_config_simple")

	assert.Equal(t, 8081, appConfiguration.Server.HTTPServer.Port)
	assert.Equal(t, 10, appConfiguration.Server.HTTPServer.ReadTimeout)
	assert.Equal(t, 11, appConfiguration.Server.HTTPServer.WriteTimeout)

	assert.Equal(t, "mongodb://localhost:27017", appConfiguration.Storage.DbConfiguration.URI)
	assert.Equal(t, "testdb", appConfiguration.Storage.DbConfiguration.Name)
	assert.Equal(t, "properties_collection", appConfiguration.Storage.DbConfiguration.PropertiesCollectionName)
}

func TestLoadEnv(t *testing.T) {
	os.Setenv("SERVER_PORT", "8081")
	os.Setenv("SERVER_READ_TIMEOUT", "10")
	os.Setenv("SERVER_WRITE_TIMEOUT", "11")

	os.Setenv("MONGO_DB_URI", "mongodb://localhost:27017")
	os.Setenv("MONGO_DB_NAME", "testdb")
	os.Setenv("MONGO_DB_PROPERTIES_COLLECTION", "properties_collection")

	appConfiguration := setupAndLoad("test_config_env")

	assert.Equal(t, 8081, appConfiguration.Server.HTTPServer.Port)
	assert.Equal(t, 10, appConfiguration.Server.HTTPServer.ReadTimeout)
	assert.Equal(t, 11, appConfiguration.Server.HTTPServer.WriteTimeout)

	assert.Equal(t, "mongodb://localhost:27017", appConfiguration.Storage.DbConfiguration.URI)
	assert.Equal(t, "testdb", appConfiguration.Storage.DbConfiguration.Name)
	assert.Equal(t, "properties_collection", appConfiguration.Storage.DbConfiguration.PropertiesCollectionName)
}

func TestLoadDefaults(t *testing.T) {
	appConfiguration := setupAndLoad("test_config_simple")

	assert.Equal(t, "text", appConfiguration.Logger.Format)
	assert.Equal(t, "info", appConfiguration.Logger.Level)
	assert.Equal(t, "./logs", appConfiguration.Logger.LogsDir)
	assert.Equal(t, "basic-go-rest-api", appConfiguration.Logger.AppLogName)
	assert.Equal(t, false, appConfiguration.Logger.AppLogConsole)
}

func setupAndLoad(name string) *AppConfiguration {
	appConfiguration := NewAppConfiguration()
	appConfiguration.Settings.configPath = "../tests/config"
	appConfiguration.Settings.configName = name

	appConfiguration.Load()

	return appConfiguration
}
