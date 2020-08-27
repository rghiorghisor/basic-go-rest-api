package logger

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/rghiorghisor/basic-go-rest-api/config"
	logrus "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func TestNewDefinitionTextFormat(t *testing.T) {
	conf := getSimpleConfiguration()

	definition := NewDefinition(conf)
	fmtr := definition.formatter
	assert.NotEqual(t, fmtr.(*prefixed.TextFormatter), nil)

	// Test case sensivity.
	conf.Format = "tExT"
	definition = NewDefinition(conf)
	fmtr = definition.formatter
	assert.NotEqual(t, fmtr.(*prefixed.TextFormatter), nil)
}

func TestNewDefinitionJSONFormat(t *testing.T) {
	conf := getSimpleConfiguration()
	conf.Format = "json"

	definition := NewDefinition(conf)
	fmtr := definition.formatter
	assert.NotEqual(t, fmtr.(*logrus.JSONFormatter), nil)

	// Test case sensivity.
	conf.Format = "jSOn"
	definition = NewDefinition(conf)
	fmtr = definition.formatter
	assert.NotEqual(t, fmtr.(*logrus.JSONFormatter), nil)

}

func TestNewDefinitionDefaultFormat(t *testing.T) {
	conf := getSimpleConfiguration()
	conf.Format = "none"

	definition := NewDefinition(conf)
	fmtr := definition.formatter
	assert.NotEqual(t, fmtr.(*prefixed.TextFormatter), nil)
}

func TestNewDefinitionLevel(t *testing.T) {
	conf := getSimpleConfiguration()

	definition := NewDefinition(conf)
	assert.Equal(t, definition.level, logrus.InfoLevel)

	// Test case sensivity.
	conf.Level = "TrAcE"
	definition = NewDefinition(conf)
	assert.Equal(t, definition.level, logrus.TraceLevel)
}

func TestNewDefinitionDefaultLevel(t *testing.T) {
	conf := getSimpleConfiguration()
	conf.Level = "zzz"

	definition := NewDefinition(conf)
	assert.Equal(t, definition.level, logrus.InfoLevel)
}

func TestNewDefinitionFile(t *testing.T) {
	conf := getSimpleConfiguration()

	definition := NewDefinition(conf)
	assert.Equal(t, definition.fileName, "logs/log-name.log")

	conf.AppLogName = "log-name-2.log"
	definition = NewDefinition(conf)
	assert.Equal(t, definition.fileName, "logs/log-name-2.log")
}

func getSimpleConfiguration() *config.LoggerConfiguration {
	return &config.LoggerConfiguration{
		Format:        "text",
		Level:         "info",
		LogsDir:       "./logs",
		AppLogName:    "log-name",
		AppLogConsole: true,
	}
}
