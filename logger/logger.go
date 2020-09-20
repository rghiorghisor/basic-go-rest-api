package logger

import (
	"io"
	"os"

	"github.com/rghiorghisor/basic-go-rest-api/config"
	"github.com/rghiorghisor/basic-go-rest-api/util"
	logrus "github.com/sirupsen/logrus"
)

// Main is the main application logger that should contain any important information.
var Main *Logger

// Access is the application access logger that should contain the requests information.
var Access *Logger

// Logger the struct containing all available loggers.
type Logger struct {
	Logger *logrus.Logger
	prefix string
}

// New generated a new logger ready to be used.
func New(conf *config.LoggersConfiguration) *Logger {
	lgrDefinition := NewDefinition(conf.MainLogger)
	lgr, _ := NewFileLogger(lgrDefinition)
	Main = lgr

	lgrDefinition = NewDefinition(conf.AccessLogger)
	lgr, _ = NewFileLogger(lgrDefinition)
	Access = lgr

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

	newLogger.Out.Write([]byte("\n"))
	lgr := &Logger{
		Logger: newLogger,
		prefix: lgrDefinition.prefix,
	}

	return lgr
}

// NewFileLogger generates a new logger that writes is a log file.
func NewFileLogger(lgrDefinition *Definition) (*Logger, error) {
	fileName := lgrDefinition.fileName

	// Create the directory in case it does not exists.
	if err := util.CreateParentFolder(fileName); err != nil {
		return nil, err
	}

	fd, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return NewLogger(lgrDefinition, fd), nil
}

// NewDummyLogger returns a dummy logger to be used, mainly, in testing.
func NewDummyLogger(w io.Writer) *Logger {
	lgrDefinition := &Definition{
		level:       logrus.InfoLevel,
		formatter:   &logrus.TextFormatter{},
		withConsole: false,
		prefix:      "main",
	}

	return NewLogger(lgrDefinition, w)
}

// Info logs the given message, along with the prefix field.
func (lgr *Logger) Info(message string) {
	lgr.withFields().
		Info(message)
}

// Infof formats the give string along with the parameters and logs it as an INFO message,
// along with the along with the prefix field.
func (lgr *Logger) Infof(format string, args ...interface{}) {
	lgr.withFields().
		Infof(format, args...)
}

// Warn logs the given message as an warning.
func (lgr *Logger) Warn(message string) {
	lgr.withFields().
		Warn(message)
}

// Errore logs the given message as an error, along with the prefix field.
func (lgr *Logger) Errore(message string) {
	lgr.withFields().
		Error(message)
}

// Error logs the given message as an error, along with the prefix field.
func (lgr *Logger) Error(message string, err error) {
	lgr.withFields().
		WithError(err).
		Error(message)
}

// Debugf formats the give string along with the parameters and logs it as a DEBUG message,
// along with the along with the prefix field.
func (lgr *Logger) Debugf(format string, args ...interface{}) {
	lgr.withFields().
		Debugf(format, args...)
}

func (lgr *Logger) withFields() *logrus.Entry {
	return lgr.Logger.WithFields(logrus.Fields{
		"prefix": lgr.prefix,
	})
}
