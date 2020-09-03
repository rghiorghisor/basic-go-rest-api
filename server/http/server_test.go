package http

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rghiorghisor/basic-go-rest-api/config"
	"github.com/rghiorghisor/basic-go-rest-api/logger"
	"github.com/rghiorghisor/basic-go-rest-api/server"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	buf := new(bytes.Buffer)
	logger.Main = logger.NewDummyLogger(buf)

	serverConfiguration := &config.HTTPServerConfiguration{
		Port: 0,
	}

	var controller server.Controller
	controller = &DummyController{}

	instance := &server.Controllers{}
	instance.HTTP = append(instance.HTTP, controller)

	srv := NewAppServer()
	srv.Setup(serverConfiguration, instance)
	go srv.Run()

	defer func() {
		p, err := os.FindProcess(os.Getpid())
		if err != nil {
			srv.httpServer.Close()
		} else {
			p.Signal(os.Interrupt)
		}

	}()

	// Wait until the server starts. If the server is not in place
	c1 := make(chan int, 1)
	go func() {
		port := ss(srv)
		c1 <- port
	}()
	address := ":"

	select {
	case res := <-c1:
		address = address + strconv.Itoa(res)
	case <-time.After(3 * time.Second):
		assert.Fail(t, "Server did not start.")
		return
	}

	testConnection(t, address)
	testResponse(t, address)
}

func testConnection(t *testing.T, address string) {
	timeout := time.Second

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		assert.Fail(t, "Cannot connect", err)
	}
	if conn != nil {
		defer conn.Close()
	}
}

func testResponse(t *testing.T, address string) {
	resp, err := http.Get("http://localhost" + address + "/api/property")
	if err != nil {
		assert.Fail(t, "Cannot connect", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		assert.Fail(t, "Cannot connect", err)
	}

	assert.Equal(t, "OK", string(body))
}

type DummyController struct {
}

// Register this controller to the provided group.
func (ctrl *DummyController) Register(routerGroup *gin.RouterGroup) {
	api := routerGroup.Group("/property")

	api.GET("", ctrl.Create)
}

func (ctrl *DummyController) Create(ctx *gin.Context) {
	ctx.String(http.StatusCreated, "OK")
}

func ss(srv *AppServer) int {
	for {
		if srv != nil && srv.listener != nil && srv.listener.Addr() != nil {
			return srv.listener.Addr().(*net.TCPAddr).Port
		}
	}
}
