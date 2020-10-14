package appserver

import (
	"log"

	"github.com/rghiorghisor/basic-go-rest-api/config"
	"github.com/rghiorghisor/basic-go-rest-api/container"
	"github.com/rghiorghisor/basic-go-rest-api/logger"
	"github.com/rghiorghisor/basic-go-rest-api/server"
	"github.com/rghiorghisor/basic-go-rest-api/server/http"
)

// AppServer is the main application server controller. It must contain in any
// version all a caller needs to start a server that exposes an application.
type AppServer struct {
	Container     *container.Container
	Configuration *config.AppConfiguration
}

// New return a basic default AppServer.
func New() *AppServer {
	return &AppServer{
		Container: container.New(),
	}
}

// LoadConfig processes the default configuration for the application server.
func (appServer *AppServer) LoadConfig() {
	// Load and validate configuration.
	appConfiguration := config.NewAppConfiguration()
	if err := appConfiguration.Load(); err != nil {
		log.Fatalf("Cannot configure server: %s", err)
	}

	appServer.Configuration = appConfiguration

	// Setup and start the logger.
	startLogger(appServer.Configuration)
}

// Start processes the settings and starts the server.
func (appServer *AppServer) Start() {
	if appServer.Configuration == nil {
		log.Fatal("[Failed to start] No configuration loaded. Please make sure AppServer.LoadConfig() is called before AppServer.Start()")
	}

	appServer.Container.Invoke(func(server *http.Server, ctls server.Controllers) {
		server.Setup(appServer.Configuration, &ctls)
		if err := server.Run(); err != nil {
			log.Fatalf("[Failed to start] %+v", err)
		}
	})
}

func startLogger(appConfiguration *config.AppConfiguration) {
	logger.New(appConfiguration.Loggers)

	logger.Main.Info("Loaded configuration.")
	logger.Main.Info(appConfiguration.Stats())
	logger.Main.Infof("Starting application %s v.%d (%s mode)...",
		appConfiguration.Application.Name,
		appConfiguration.Application.Version,
		appConfiguration.Environment.Name)
}
