package config

// NewDefault creates an AppConfiguration object that contains all the default
// configuration values.
func NewDefault() *AppConfiguration {
	return &AppConfiguration{
		Settings:    newDefaultSettings(),
		Application: newDefaultApplicationSettings(),
		Loggers:     newDefaultLoggersConfiguration(),
		Storage:     newDefaultStorageConfiguration(),
		Server:      newDefaultServerConfiguration(),
	}
}

func newDefaultSettings() *ConfigurationSettings {
	return &ConfigurationSettings{
		configPath: "./config",
		configName: "config",
	}
}

func newDefaultApplicationSettings() *ApplicationSettings {
	return &ApplicationSettings{
		Name:        "Basic Go REST API",
		Version:     1,
		ContextPath: "/api/v1",
	}
}

func newDefaultLoggersConfiguration() *LoggersConfiguration {
	return &LoggersConfiguration{
		MainLogger:   newDefaultMainLoggerConfiguration(),
		AccessLogger: newDefaultAccessLoggerConfiguration(),
	}
}

func newDefaultMainLoggerConfiguration() *LoggerConfiguration {
	return &LoggerConfiguration{
		Format:      "text",
		Level:       "info",
		LogsDir:     "./logs",
		FileName:    "basic-go-rest-api",
		WithConsole: false,
		Prefix:      "main",
	}
}

func newDefaultAccessLoggerConfiguration() *LoggerConfiguration {
	return &LoggerConfiguration{
		Format:      "text",
		Level:       "info",
		LogsDir:     "./logs",
		FileName:    "access",
		WithConsole: false,
		Prefix:      "access",
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
