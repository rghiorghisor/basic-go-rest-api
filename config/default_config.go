package config

// NewDefault creates an AppConfiguration object that contains all the default
// configuration values.
func NewDefault() *AppConfiguration {
	return &AppConfiguration{
		Settings: newDefaultSettings(),
		Logger:   newDefaultLoggerConfiguration(),
		Server:   newDefaultServerConfiguration(),
	}
}

func newDefaultSettings() *ConfigurationSettings {
	return &ConfigurationSettings{"./config", "config"}
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
