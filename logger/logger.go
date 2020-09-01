package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/rghiorghisor/basic-go-rest-api/config"
	logrus "github.com/sirupsen/logrus"
)

// Main is the main application logger that should contain any important information.
var Main *Logger

// Logger the struct containing all available loggers.
type Logger struct {
	Logger *logrus.Logger
	prefix string
}

// New generated a new logger ready to be used.
func New(conf *config.LoggerConfiguration) *Logger {
	lgrDefinition := NewDefinition(conf)

	lgr, _ := NewFileLogger(lgrDefinition)

	return lgr
}

// NewLogger creates a new logger based on the provided parameters.
func NewLogger(lgrDefinition *Definition, w io.Writer) *Logger {
	newLogger := logrus.New()

	newLogger.Level = lgrDefinition.level
	newLogger.Formatter = lgrDefinition.formatter

	if lgrDefinition.withConsole && w != os.Stdout {
		mw := io.MultiWriter(w, os.Stdout)
		newLogger.Out = mw
	} else {
		newLogger.Out = w
	}

	lgr := &Logger{
		Logger: newLogger,
	}

	return lgr
}

// NewFileLogger generates a new logger that writes is a log file.
func NewFileLogger(lgrDefinition *Definition) (*Logger, error) {
	fileName := lgrDefinition.fileName

	// Create the directory in case it does not exists.
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}

	fd, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	Main = NewLogger(lgrDefinition, fd)
	Main.prefix = "main"

	return Main, nil
}

// Info logs the given message, along with the prefix field.
func (lgr *Logger) Info(message string) {
	lgr.Logger.WithFields(logrus.Fields{
		"prefix": lgr.prefix,
	}).Info(message)
}

// Infof formats the give string along with the parameters and logs it, along
// with the along with the prefix field.
func (lgr *Logger) Infof(format string, args ...interface{}) {
	lgr.Logger.WithFields(logrus.Fields{
		"prefix": lgr.prefix,
	}).Infof(format, args...)
}

// Error logs the given message as an error, along with the prefix field.
func (lgr *Logger) Error(message string, err error) {
	lgr.Logger.WithFields(logrus.Fields{
		"prefix": lgr.prefix,
	}).WithError(err).Error(message)
}
