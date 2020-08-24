package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rghiorghisor/basic-go-rest-api/config"
	phttp "github.com/rghiorghisor/basic-go-rest-api/property/gateway/http"
	shttp "github.com/rghiorghisor/basic-go-rest-api/server/http"
)

// AppServer structure that encapsulates all related data.
type AppServer struct {
	httpServer *http.Server
}

// NewAppServer creates a new bare-boned application server.
func NewAppServer() *AppServer {
	return &AppServer{}
}

// Setup prepares the server but does not starts it.
func (server *AppServer) Setup(serverConfiguration *config.HTTPServerConfiguration, services *Services) {
	// Initialize the gin router.
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		shttp.JSONAppErrorHandler(),
	)

	setupEndpoints(services, router)
	setupServer(server, router, serverConfiguration)
}

// Run starts the application server based on configuration settings.
func (server *AppServer) Run() {

	go func() {
		if err := server.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	server.httpServer.Shutdown(ctx)
}

func setupEndpoints(services *Services, router *gin.Engine) {
	api := router.Group("/api")

	// Register endpoints to api.
	phttp.RegisterEndpoints(api, services.PropertyService)
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
