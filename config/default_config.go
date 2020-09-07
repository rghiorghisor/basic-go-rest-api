package config

// NewDefault creates an AppConfiguration object that contains all the default
// configuration values.
func NewDefault() *AppConfiguration {
	return &AppConfiguration{
		Settings: newDefaultSettings(),
		Logger:   newDefaultLoggerConfiguration(),
		Storage:  newDefaultStorageConfiguration(),
		Server:   newDefaultServerConfiguration(),
	}
}

func newDefaultSettings() *ConfigurationSettings {
	return &ConfigurationSettings{
		configPath: "./config",
		configName: "config",
	}
}

func newDefaultLoggerConfiguration() *LoggerConfiguration {
	return &LoggerConfiguration{
		Format:        "text",
		Level:         "info",
		LogsDir:       "./logs",
		AppLogName:    "basic-go-rest-api",
		AppLogConsole: false,
	}
}

func newDefaultServerConfiguration() *ServerConfiguration {
	return &ServerConfiguration{
		HTTPServer: newDefaultHTTPServerConfiguration(),
	}
}

func newDefaultHTTPServerConfiguration() *HTTPServerConfiguration {
	return &HTTPServerConfiguration{
		Port:         8080,
		ReadTimeout:  10,
		WriteTimeout: 10,
	}
}

func newDefaultStorageConfiguration() *StorageConfiguration {
	return &StorageConfiguration{
		BoltDbConfiguration: newDefaultBoltDbConfiguration(),
	}
}

func newDefaultBoltDbConfiguration() *BoltDbConfiguration {
	return &BoltDbConfiguration{
		Name: "local-storage/boltdb",
	}
}
