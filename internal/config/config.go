package config

import ()

type (
	IConfig interface {
		Get(key string) any
	}
	// AppConfig struct {
	// 	Environment     string
	// 	LogLevel        string        `envconfig:"LOG_LEVEL" default:"DEBUG"`
	// 	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"60s"`
	// }
)
