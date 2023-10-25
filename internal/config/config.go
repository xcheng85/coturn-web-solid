package config

import (
	// "fmt"
	// "os"
	// "time"
)

type (
	IConfig interface {
		Get(key string) any
	}
	// AppConfig struct {
	// 	Environment     string
	// 	LogLevel        string        `envconfig:"LOG_LEVEL" default:"DEBUG"`
	// 	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"60s"`
	// }

	// WebConfig struct {
	// 	Host string `default:"127.0.0.1"` // 0.0.0.0 has wsl issues, 127.0.0.1 has linux issue
	// 	Port string `default:":8080"`
	// }
	// Config struct {
	// 	// AppConfig // type embedding
	// 	// WebConfig
	// }
)

// func (c WebConfig) Address() string {
// 	return fmt.Sprintf("%s%s", c.Host, c.Port)
// }

