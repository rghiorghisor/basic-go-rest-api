package storage

import (
	"os"
	"testing"

	"github.com/rghiorghisor/basic-go-rest-api/config"
	"github.com/rghiorghisor/basic-go-rest-api/logger"
	"github.com/stretchr/testify/mock"
)

var tContext *testContext

func TestSetup(t *testing.T) {
	tContext = new(testContext)
	tContext.mock1 = new(factoryMock)
	tContext.mock2 = new(factoryMock)
	tContext.mock3 = new(factoryMock)

	logger.Main = logger.NewDummyLogger(os.Stdout)

	storage := New()
	storage.factories = []func() factory{newMock1, newMock2}
	storage.defaultFactory = newMock3

	cfg := &config.StorageConfiguration{Type: "local"}
	tContext.mock1.On("id").Return("local")
	tContext.mock1.On("init", storage, cfg).Return(nil)
	tContext.mock2.On("id").Return("mongo")
	tContext.mock2.On("init", storage, cfg).Return(nil)
	tContext.mock3.On("id").Return("local")
	tContext.mock3.On("init", storage, cfg).Return(nil)

	storage.SetupStorage(cfg)

	tContext.mock1.AssertNumberOfCalls(t, "init", 1)
	tContext.mock2.AssertNumberOfCalls(t, "init", 0)
	tContext.mock3.AssertNumberOfCalls(t, "init", 0)

	cfg.Type = "mongo"
	storage.SetupStorage(cfg)
	tContext.mock1.AssertNumberOfCalls(t, "init", 1)
	tContext.mock2.AssertNumberOfCalls(t, "init", 1)
	tContext.mock3.AssertNumberOfCalls(t, "init", 0)

	cfg.Type = "none"
	storage.SetupStorage(cfg)
	tContext.mock1.AssertNumberOfCalls(t, "init", 1)
	tContext.mock2.AssertNumberOfCalls(t, "init", 1)
	tContext.mock3.AssertNumberOfCalls(t, "init", 1)

}

type testContext struct {
	mock1 *factoryMock
	mock2 *factoryMock
	mock3 *factoryMock
}

type factoryMock struct {
	mock.Mock
}

func newMock1() factory {
	return tContext.mock1
}

func newMock2() factory {
	return tContext.mock2
}

func newMock3() factory {
	return tContext.mock3
}

func (m *factoryMock) id() string {
	args := m.Called()

	return args.String(0)
}

func (m *factoryMock) init(storage *Storage, config *config.StorageConfiguration) error {
	args := m.Called(storage, config)

	return args.Error(0)
}
