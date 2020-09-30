package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/rghiorghisor/basic-go-rest-api/util"
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

var flags = createFlags()
var flagsInit = initFlags()

type appFlag struct {
	command      string
	defaultValue string
	description  string
	ref          string
	validator    func(*appFlag, *AppConfiguration, string) error
}

// AppConfiguration holds the current application configuration.
type AppConfiguration struct {
	Environment *Environment `yaml:"none"`
	Settings    *ConfigurationSettings
	Loggers     *LoggersConfiguration `yaml:"logger"`
	Server      *ServerConfiguration  `yaml:"server"`
	Storage     *StorageConfiguration `yaml:"storage"`
	stats       *stats
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

// LoggersConfiguration contains the settings of all available loggers.
type LoggersConfiguration struct {
	MainLogger   *LoggerConfiguration `yaml:"main"`
	AccessLogger *LoggerConfiguration `yaml:"access"`
}

// LoggerConfiguration contains all available configurable logging settings.
type LoggerConfiguration struct {
	Format      string `yaml:"format"`
	Level       string `yaml:"level"`
	LogsDir     string `yaml:"dir"`
	FileName    string `yaml:"file-name"`
	WithConsole bool   `yaml:"with-console"`
	Prefix      string `yaml:"prefix"`
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
	Type                string                `yaml:"type"`
	BoltDbConfiguration *BoltDbConfiguration  `yaml:"bolt"`
	DbConfiguration     *MongoDbConfiguration `yaml:"mongo"`
}

// BoltDbConfiguration hols settings referring to BoltDB as storage option.
type BoltDbConfiguration struct {
	Name string `yaml:"name"`
}

// MongoDbConfiguration hols settings referring to MongoDb as storage option.
type MongoDbConfiguration struct {
	URI  string `yaml:"uri"`
	Name string `yaml:"name"`
}

type stats struct {
	loaded         bool
	loadedFromDir  string
	loadedFromFile string
	foundInEnv     []interface{}
	notFoundInEnv  []interface{}
}

// NewAppConfiguration creates a new application configuration instance.
func NewAppConfiguration() *AppConfiguration {
	new := NewDefault()

	defaultEnvironment = environments[DefaultEnvCode]
	new.Environment = defaultEnvironment

	new.stats = &stats{}

	return new
}

// Load the application configuration.
//
// The loading process takes into consideration the appConfiguration.Settings
// to determine what and how to load.
func (appConfiguration *AppConfiguration) Load() error {
	if err := appConfiguration.processFlags(); err != nil {
		return err
	}

	configPath := appConfiguration.Settings.configPath
	configName := appConfiguration.Settings.configName

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)

	viper.SetConfigType(configType)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	envLoader := newEnvValueLoader()
	err := viper.Unmarshal(&appConfiguration, func(config *mapstructure.DecoderConfig) {
		config.TagName = configType
	}, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
				found, err := envLoader.load(appConfiguration, f, t, data)

				return found, err
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

	appConfiguration.stats.loaded = true
	appConfiguration.stats.loadedFromDir = configPath
	appConfiguration.stats.loadedFromFile = configName

	return err
}

func createFlags() []*appFlag {
	return []*appFlag{
		{command: "env", defaultValue: "dev", description: "Environment {prod|production|dev|development}", validator: validateEnvFlag},
		{command: "c", defaultValue: "config/config.yml", description: "Path and name of the configuration file", validator: validateConfigFlag},
	}
}

func initFlags() error {
	for _, f := range flags {
		flag.StringVar(&f.ref, f.command, f.defaultValue, f.description)
	}

	return nil
}

func (appConfiguration *AppConfiguration) processFlags() error {
	flag.Parse()

	for _, f := range flags {
		if !isFlagPassed(f.command) {
			continue
		}

		if err := f.validator(f, appConfiguration, f.ref); err != nil {
			return err
		}
	}

	return nil
}

func isFlagPassed(name string) bool {
	found := false

	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})

	return found
}

func (appConfiguration *AppConfiguration) processEnvironment() {
	envString := viper.GetString("environment")

	appConfiguration.Environment = appConfiguration.validateRawEnvironment(envString)
	if appConfiguration.Environment != nil {
		return
	}

	appConfiguration.Environment = defaultEnvironment
	fmt.Printf("[WARNING] Unknown environment='%s'. Using default '%s'.\r\n", envString, defaultEnvironment.Name)
}

func (appConfiguration *AppConfiguration) validateRawEnvironment(envString string) *Environment {
	for _, s := range environments {
		for _, acc := range s.accepted {
			if strings.EqualFold(envString, acc) {
				return s
			}
		}
	}

	return nil
}

func validateEnvFlag(flg *appFlag, appConfiguration *AppConfiguration, envString string) error {
	e := appConfiguration.validateRawEnvironment(envString)
	if e == nil {
		return newError(flg, envString)
	}

	appConfiguration.Settings.environment = envString
	return nil
}

func validateConfigFlag(flg *appFlag, appConfiguration *AppConfiguration, configFile string) error {
	dir, file := util.SplitFilepath(configFile)

	appConfiguration.Settings.configPath = dir
	appConfiguration.Settings.configName = file

	return nil
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

func (envValueLoader *EnvValueLoader) load(appConfiguration *AppConfiguration, f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
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

	if value == "" {
		appConfiguration.stats.notFoundInEnv = append(appConfiguration.stats.notFoundInEnv, valueString)
	} else {
		appConfiguration.stats.foundInEnv = append(appConfiguration.stats.foundInEnv, valueString)
	}

	return value, nil
}

func newError(flg *appFlag, value string) error {
	s := "Invalid values for flag %s ('%s'). %s description: '%s'"
	s = fmt.Sprintf(s, flg.command, value, flg.command, flg.description)

	return errors.New(s)
}

// Stats retrieves the configuration loading statistics.
func (appConfiguration *AppConfiguration) Stats() string {
	return fmt.Sprintf("%+v", *appConfiguration.stats)
}
