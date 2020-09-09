package logger

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rghiorghisor/basic-go-rest-api/config"
	logrus "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	defaultLogExtension    = ".log"
	defaultTimestampFormat = "2006-01-02 15:04:05.000"
)

// Definition contains all settings that must be used when defining a logger.
// It is the responsibility of this structure to normalize the values and the
// logger creation it self must take the values as they are.
type Definition struct {
	fileName    string
	formatter   logrus.Formatter
	level       logrus.Level
	withConsole bool
	prefix      string
}

// NewDefinition loads the values from the provided configuration, processes them
// and retrieves a new definition struct.
func NewDefinition(conf *config.LoggerConfiguration) *Definition {
	fileName := getFileName(conf)
	formatter := getFormatter(conf)
	level := getLevel(conf)

	return &Definition{
		fileName:    fileName,
		formatter:   formatter,
		level:       level,
		withConsole: conf.WithConsole,
		prefix:      conf.Prefix,
	}
}

func getFileName(conf *config.LoggerConfiguration) string {
	dirString := conf.LogsDir
	fileName := conf.FileName

	// Add file extension in case it is not already there.
	if !strings.HasSuffix(fileName, defaultLogExtension) {
		fileName = fileName + defaultLogExtension
	}

	f := filepath.Join(dirString, fileName)

	return f
}

func getFormatter(conf *config.LoggerConfiguration) logrus.Formatter {
	formatterString := conf.Format

	if strings.EqualFold(formatterString, "json") {
		return &logrus.JSONFormatter{
			TimestampFormat: defaultTimestampFormat,
		}
	}

	textFormatter := &prefixed.TextFormatter{
		TimestampFormat: defaultTimestampFormat,
		DisableColors:   false,
		ForceColors:     true,
		FullTimestamp:   true,
		ForceFormatting: true,
	}

	if strings.EqualFold(formatterString, "text") {
		return textFormatter
	}

	err := fmt.Errorf("[ERROR] logger: Invalid log formatter ('%s'). Using default '%s'", formatterString, "prefixed.TextFormatter")
	fmt.Println(err)
	return textFormatter
}

func getLevel(conf *config.LoggerConfiguration) logrus.Level {
	levelString := conf.Level

	level, err := logrus.ParseLevel(levelString)
	if err != nil {

		err := fmt.Errorf("[ERROR] logger: Invalid log format ('%s'). Using default '%v'", levelString, logrus.InfoLevel)
		fmt.Println(err)
		return logrus.InfoLevel
	}

	return level
}
