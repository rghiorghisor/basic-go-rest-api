package logger

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	logrus "github.com/sirupsen/logrus"
)

func TestInfo(t *testing.T) {
	definition := &Definition{
		level:     logrus.InfoLevel,
		formatter: &logrus.TextFormatter{},
		prefix:    "main",
	}

	buf := new(bytes.Buffer)
	lgr := NewLogger(definition, buf)

	lgr.Info("test message")

	actual := buf.String()

	assert.Equal(t, strings.Contains(actual, "level=info"), true)
	assert.Equal(t, strings.Contains(actual, "prefix=main"), true)
	assert.Equal(t, strings.Contains(actual, "msg=\"test message\""), true)
}

func TestInfof(t *testing.T) {
	definition := &Definition{
		level:     logrus.DebugLevel,
		formatter: &logrus.TextFormatter{},
		prefix:    "main",
	}

	buf := new(bytes.Buffer)
	lgr := NewLogger(definition, buf)

	lgr.Infof("test message no %d", 1)

	actual := buf.String()
	assert.Equal(t, strings.Contains(actual, "level=info"), true)
	assert.Equal(t, strings.Contains(actual, "prefix=main"), true)
	assert.Equal(t, strings.Contains(actual, "msg=\"test message no 1\""), true)
}

func Info1(t *testing.T) {
	definition := &Definition{
		level:       logrus.DebugLevel,
		formatter:   &logrus.TextFormatter{},
		withConsole: true,
		prefix:      "main",
	}

	buf := new(bytes.Buffer)
	lgr := NewLogger(definition, buf)

	lgr.Info("test message")

	actual := captureOutput(lgr, func() {
		lgr.Info("test message")
	})

	actual = "\n" + actual
	assert.Equal(t, actual, buf.String())
}

func TestInfoFile(t *testing.T) {
	testFileLocation := "../tests/tmp/tmp.log"
	testDirLocation := "../tests/tmp"

	definition := &Definition{
		level:       logrus.DebugLevel,
		formatter:   &logrus.TextFormatter{},
		withConsole: true,
		fileName:    testFileLocation,
	}

	NewFileLogger(definition)

	info, err := os.Stat(testFileLocation)

	assert.Equal(t, info.IsDir(), false)
	assert.Equal(t, os.IsNotExist(err), false)

	defer func() {
		os.Remove(testFileLocation)
		os.Remove(testDirLocation)
	}()
}

func captureOutput(lgr *Logger, f func()) string {
	old := lgr.Logger.Out
	r, w, _ := os.Pipe()
	lgr.Logger.Out = w

	f()

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	lgr.Logger.Out = old
	out := <-outC

	return out
}
