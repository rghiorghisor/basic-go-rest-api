package logger

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	logrus "github.com/sirupsen/logrus"
)

func TestInfo(t *testing.T) {
	buf, lgr := simpleSetup(logrus.InfoLevel)

	lgr.Info("test message")

	actual := buf.String()

	assert.Equal(t, true, strings.Contains(actual, "level=info"))
	assert.Equal(t, true, strings.Contains(actual, "prefix=main"))
	assert.Equal(t, true, strings.Contains(actual, "msg=\"test message\""))
}

func TestInfof(t *testing.T) {
	buf, lgr := simpleSetup(logrus.InfoLevel)

	lgr.Infof("test message no %d", 1)

	actual := buf.String()
	assert.Equal(t, true, strings.Contains(actual, "level=info"))
	assert.Equal(t, true, strings.Contains(actual, "prefix=main"))
	assert.Equal(t, true, strings.Contains(actual, "msg=\"test message no 1\""))
}

func Info1(t *testing.T) {
	buf, lgr := simpleSetup(logrus.DebugLevel)
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

func TestWarn(t *testing.T) {
	buf, lgr := simpleSetup(logrus.InfoLevel)

	lgr.Warn("test message")

	actual := buf.String()
	assert.Equal(t, true, strings.Contains(actual, "level=warn"))
	assert.Equal(t, true, strings.Contains(actual, "prefix=main"))
	assert.Equal(t, true, strings.Contains(actual, "msg=\"test message\""))
}

func TestDebugf(t *testing.T) {
	buf, lgr := simpleSetup(logrus.DebugLevel)

	lgr.Debugf("test message no %d", 1)

	actual := buf.String()
	assert.Equal(t, true, strings.Contains(actual, "level=debug"))
	assert.Equal(t, true, strings.Contains(actual, "prefix=main"))
	assert.Equal(t, true, strings.Contains(actual, "msg=\"test message no 1\""))
}

func TestErrore(t *testing.T) {
	buf, lgr := simpleSetup(logrus.DebugLevel)

	lgr.Errore("test message")

	actual := buf.String()
	assert.Equal(t, true, strings.Contains(actual, "level=error"))
	assert.Equal(t, true, strings.Contains(actual, "prefix=main"))
	assert.Equal(t, true, strings.Contains(actual, "msg=\"test message\""))
}

func TestError(t *testing.T) {
	buf, lgr := simpleSetup(logrus.DebugLevel)

	lgr.Error("test message", errors.New("unexpected"))

	actual := buf.String()

	fmt.Println(actual)

	assert.Equal(t, true, strings.Contains(actual, "level=error"))
	assert.Equal(t, true, strings.Contains(actual, "prefix=main"))
	assert.Equal(t, true, strings.Contains(actual, "msg=\"test message\""))
	assert.Equal(t, true, strings.Contains(actual, "error=unexpected"))
}

func simpleSetup(level logrus.Level) (*bytes.Buffer, *Logger) {
	definition := &Definition{
		level:     level,
		formatter: &logrus.TextFormatter{},
		prefix:    "main",
	}

	buf := new(bytes.Buffer)
	lgr := NewLogger(definition, buf)

	return buf, lgr
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
