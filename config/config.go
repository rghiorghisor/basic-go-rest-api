package config

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// configType defines the type of configuration. As of this moment this is a
// constant as AppConfiguration supports a single configuration file type.
const configType = "yaml"

const (
	developCode    = 0
	productionCode = 1
)

var environments = []*Environment{
	{code: developCode, Name: "development", accepted: []string{"dev", "development"}},
	{code: productionCode, Name: "production", accepted: []string{"prod", "production"}},
}
var defaultEnvironment = environments[0]

// AppConfiguration holds the current application configuration.
type AppConfiguration struct {
	Environment *Environment
	Settings    *ConfigurationSettings
	Logger      *LoggerConfiguration  `yaml:"logger"`
	Server      *ServerConfiguration  `yaml:"server"`
	Storage     *StorageConfiguration `yaml:"storage"`
}

// Environment represents the settings of a particular application working mode.
type Environment struct {
	code     int
	Name     string
	accepted []string
}

// ConfigurationSettings contains some metadata information needed when loading
// the configuration.
type ConfigurationSettings struct {
	configPath  string
	configName  string
	environment string
}

// LoggerConfiguration contains all available configurable logging settings.
type LoggerConfiguration struct {
	Format        string `yaml:"format"`
	Level         string `yaml:"level"`
	LogsDir       string `yaml:"dir"`
	AppLogName    string `yaml:"application-log-file-name"`
	AppLogConsole bool   `yaml:"application-log-console"`
}

// ServerConfiguration holds any settings regarding the application's server.
type ServerConfiguration struct {
	HTTPServer *HTTPServerConfiguration `yaml:"http"`
}

// HTTPServerConfiguration holds settings of the HTTP specific server.
type HTTPServerConfiguration struct {
	Port         int `yaml:"port"`
	ReadTimeout  int `yaml:"read-timeout"`
	WriteTimeout int `yaml:"write-timeout"`
}

// StorageConfiguration holds any settings regarding the application's storage options.
type StorageConfiguration struct {
	DbConfiguration *MongoDbConfiguration `yaml:"mongo"`
}

// MongoDbConfiguration hols settings referring to MongoDb as storage option.
type MongoDbConfiguration struct {
	URI                      string `yaml:"uri"`
	Name                     string `yaml:"name"`
	PropertiesCollectionName string `yaml:"properties-collection"`
}

// NewAppConfiguration creates a new application configuration instance.
func NewAppConfiguration() *AppConfiguration {
	return NewDefault()
}

// Load the application configuration.
//
// The loading process takes into consideration the appConfiguration.Settings
// to determine what and how to load.
func (appConfiguration *AppConfiguration) Load() error {
	configPath := appConfiguration.Settings.configPath
	configName := appConfiguration.Settings.configName

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)

	viper.SetConfigType(configType)

	viper.ReadInConfig()

	envLoader := newEnvValueLoader()
	err := viper.Unmarshal(&appConfiguration, func(config *mapstructure.DecoderConfig) {
		config.TagName = configType
	}, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
				return envLoader.load(f, t, data)
			},
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	))

	environmentString := appConfiguration.Settings.environment
	if environmentString != "" {
		viper.Set("environment", environmentString)
	}
	appConfiguration.processEnvironment()

	return err
}

func (appConfiguration *AppConfiguration) processEnvironment() {
	envString := viper.GetString("environment")

	for _, s := range environments {
		for _, acc := range s.accepted {
			if strings.EqualFold(envString, acc) {
				appConfiguration.Environment = s

				return
			}
		}
	}

	fmt.Printf("[WARNING] Unknown environment='%s'. Using default '%s'.\r\n", envString, defaultEnvironment.Name)
	appConfiguration.Environment = defaultEnvironment
}

// IsProduction retrieves "true" if the current configured environment is "production",
// otherwise returns "false".
func (appConfiguration *AppConfiguration) IsProduction() bool {
	return appConfiguration.Environment.code == productionCode
}

// IsDevelopment retrieves "true" if the current configured environment is "development",
// otherwise returns "false".
func (appConfiguration *AppConfiguration) IsDevelopment() bool {
	return appConfiguration.Environment.code == developCode
}

// EnvValueLoader handles loading environment variables.
type EnvValueLoader struct {
	expression *regexp.Regexp
	reference  map[string]interface{}
}

// newEnvValueLoader creates a new instance that can be used for loading env variables.
func newEnvValueLoader() *EnvValueLoader {
	return &EnvValueLoader{expression: regexp.MustCompile(`^\$.*`)}
}

func (envValueLoader *EnvValueLoader) load(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f == reflect.TypeOf(envValueLoader.reference) {
		// If it is not a value, but a struct ignore further processing.
		return data, nil
	}

	valueString := fmt.Sprintf("%v", data)
	if !envValueLoader.expression.MatchString(valueString) {
		// If the value is not a ENV variable ignore further processing.
		return data, nil
	}

	varName := valueString[1:]
	value := os.Getenv(varName)

	return value, nil
}
