package http

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rghiorghisor/basic-go-rest-api/config"
	"github.com/rghiorghisor/basic-go-rest-api/logger"
	"github.com/rghiorghisor/basic-go-rest-api/server"
)

// AppServer structure that encapsulates all related data.
type AppServer struct {
	httpServer *http.Server
	listener   net.Listener
}

// NewAppServer creates a new bare-boned application server.
func NewAppServer() *AppServer {
	return &AppServer{}
}

// Setup prepares the server but does not starts it.
func (server *AppServer) Setup(serverConfiguration *config.HTTPServerConfiguration, controllers *server.Controllers) {
	// Initialize the gin router.
	router := gin.Default()

	router.Use(
		AccessLogger(),
		gin.Recovery(),
		gin.Logger(),
		JSONAppErrorHandler(),
	)

	setupEndpoints(controllers, router)
	setupServer(server, router, serverConfiguration)
}

// Run starts the application server based on configuration settings.
func (server *AppServer) Run() error {
	ch := make(chan error)
	go func() {
		listener, err := net.Listen("tcp", server.httpServer.Addr)
		if err != nil {
			ch <- err
			return
		}

		server.listener = listener
		logger.Main.Infof("Starting HTTP server, listening on '%v'", server.httpServer.Addr)
		if err := server.httpServer.Serve(listener); err != nil {
			ch <- err
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	select {
	case res := <-ch:
		return res
	case <-quit:
	}

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	logger.Main.Info("Shuting down...")
	if err := server.httpServer.Shutdown(ctx); err != nil {
		logger.Main.Error("Failed to shutdown server", err)
	}

	logger.Main.Info("Server shutdown.")

	return nil
}

func setupEndpoints(controllers *server.Controllers, router *gin.Engine) {
	base := router.Group("")
	healthcheck := NewHealthcheckController()
	healthcheck.Register(base)

	api := router.Group("/api")
	for _, c := range controllers.HTTP {
		c.Register(api)
	}
}

func setupServer(server *AppServer, router *gin.Engine, serverConfiguration *config.HTTPServerConfiguration) {
	address := ":" + strconv.Itoa(serverConfiguration.Port)
	readTimeout := time.Duration(serverConfiguration.ReadTimeout) * time.Second
	writeTimeout := time.Duration(serverConfiguration.WriteTimeout) * time.Second

	server.httpServer = &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}
}
