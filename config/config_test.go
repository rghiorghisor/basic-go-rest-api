package config

import (
	"errors"
	"flag"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestLoad(t *testing.T) {
	appConfiguration, _ := setupAndLoad("test_config_simple")

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
	os.Setenv("MONGO_DB_PROPERTIES_COLLECTION", "properties_collection")

	appConfiguration, _ := setupAndLoad("test_config_env")

	assert.Equal(t, 8081, appConfiguration.Server.HTTPServer.Port)
	assert.Equal(t, 10, appConfiguration.Server.HTTPServer.ReadTimeout)
	assert.Equal(t, 11, appConfiguration.Server.HTTPServer.WriteTimeout)

	assert.Equal(t, "mongodb://localhost:27017", appConfiguration.Storage.DbConfiguration.URI)
	assert.Equal(t, "properties_collection", appConfiguration.Storage.DbConfiguration.PropertiesCollectionName)
	assert.Equal(t, "{loaded:true "+
		"loadedFromDir:../tests/config "+
		"loadedFromFile:test_config_env "+
		"foundInEnv:[$SERVER_PORT $SERVER_READ_TIMEOUT $SERVER_WRITE_TIMEOUT $MONGO_DB_URI $MONGO_DB_PROPERTIES_COLLECTION] "+
		"notFoundInEnv:[$MONGO_DB_NAME]}", appConfiguration.Stats())
}

func TestLoadDefaults(t *testing.T) {
	appConfiguration, _ := setupAndLoad("test_config_simple")

	assert.Equal(t, "text", appConfiguration.Logger.Format)
	assert.Equal(t, "info", appConfiguration.Logger.Level)
	assert.Equal(t, "./logs", appConfiguration.Logger.LogsDir)
	assert.Equal(t, "basic-go-rest-api", appConfiguration.Logger.AppLogName)
	assert.Equal(t, false, appConfiguration.Logger.AppLogConsole)

	assert.Equal(t, developCode, appConfiguration.Environment.code)
}

func TestEnvironmentParam(t *testing.T) {
	appConfiguration, _ := setupAndLoad("test_config_env_prod")

	assert.Equal(t, productionCode, appConfiguration.Environment.code)
}

func TestLoadEnvironment(t *testing.T) {
	appConfiguration, _ := setupAndLoad("test_config_simple")
	assert.Equal(t, true, appConfiguration.IsDevelopment())
	assert.Equal(t, false, appConfiguration.IsProduction())

	appConfiguration, _ = setupAndLoadWithEnv("test_config_simple", "production")
	assert.Equal(t, false, appConfiguration.IsDevelopment())
	assert.Equal(t, true, appConfiguration.IsProduction())

	appConfiguration, _ = setupAndLoadWithEnv("test_config_simple", "prod")
	assert.Equal(t, false, appConfiguration.IsDevelopment())
	assert.Equal(t, true, appConfiguration.IsProduction())

	appConfiguration, _ = setupAndLoadWithEnv("test_config_simple", "dev")
	assert.Equal(t, true, appConfiguration.IsDevelopment())
	assert.Equal(t, false, appConfiguration.IsProduction())

	appConfiguration, _ = setupAndLoadWithEnv("test_config_simple", "development")
	assert.Equal(t, true, appConfiguration.IsDevelopment())
	assert.Equal(t, false, appConfiguration.IsProduction())
}

func TestEnvironmentFlag(t *testing.T) {
	flag.Set("env", "prod")
	appConfiguration, _ := setupAndLoad("test_config_env_dev")
	assert.Equal(t, productionCode, appConfiguration.Environment.code)

	flag.Set("env", "prod")
	appConfiguration, _ = setupAndLoad("test_config_env_prod")
	assert.Equal(t, productionCode, appConfiguration.Environment.code)

	flag.Set("env", "dev")
	appConfiguration, _ = setupAndLoad("test_config_env_prod")
	assert.Equal(t, developCode, appConfiguration.Environment.code)

	flag.Set("env", "dev")
	appConfiguration, _ = setupAndLoad("test_config_env_dev")
	assert.Equal(t, developCode, appConfiguration.Environment.code)

	flag.Set("env", "production")
	appConfiguration, _ = setupAndLoad("test_config_env_dev")
	assert.Equal(t, productionCode, appConfiguration.Environment.code)

	flag.Set("env", "production")
	appConfiguration, _ = setupAndLoad("test_config_env_prod")
	assert.Equal(t, productionCode, appConfiguration.Environment.code)

	flag.Set("env", "development")
	appConfiguration, _ = setupAndLoad("test_config_env_prod")
	assert.Equal(t, developCode, appConfiguration.Environment.code)

	flag.Set("env", "development")
	appConfiguration, _ = setupAndLoad("test_config_env_dev")
	assert.Equal(t, developCode, appConfiguration.Environment.code)
}

func TestEnvironmentFlagInvalid(t *testing.T) {
	flag.Set("env", "aaa")
	_, err := setupAndLoad("test_config_env_dev")
	assert.Equal(t, errors.New("Invalid values for flag env ('aaa'). env description: 'Environment {prod|production|dev|development}'"), err)
}

func setupAndLoad(name string) (*AppConfiguration, error) {
	return setupAndLoadWithEnv(name, "")
}

func setupAndLoadWithEnv(name string, env string) (*AppConfiguration, error) {
	appConfiguration := NewAppConfiguration()
	appConfiguration.Settings.configPath = "../tests/config"
	appConfiguration.Settings.configName = name
	appConfiguration.Settings.environment = env

	if err := appConfiguration.Load(); err != nil {
		return nil, err
	}

	return appConfiguration, nil

}
