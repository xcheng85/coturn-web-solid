package config

import (
	"path/filepath" // go 1.21.3+
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ViperConfig struct {
}

func (s *ViperConfig) Get(key string) any {
	return viper.Get(key)
}

// interface and implementation
func NewViperConfig(paths []string, logger *zap.Logger) (cfg IConfig, err error) {
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("/config/")
	viper.SetDefault("port", 8080)
	err = viper.ReadInConfig()
	viper.WatchConfig()
	if err != nil {
		logger.Sugar().Panicw("fatal error config file: %w", err)
	}
	for _, p := range paths {
		dir, file := filepath.Split(p)
		splitR := strings.Split(file, ".")
		var filetype, filename string
		if len(splitR) == 1 {
			filetype = "json"
		} else {
			filetype = splitR[1]
		}
		filename = splitR[0]
		v := viper.New()
		v.SetConfigName(filename)
		v.SetConfigType(filetype)
		v.AddConfigPath(dir)
		err := v.ReadInConfig()
		if err != nil {
			logger.Sugar().Panicw("fatal error config file: %w", err)
		}
		viper.MergeConfigMap(v.AllSettings())
	}
	logger.Sugar().Info(viper.AllKeys());
	logger.Sugar().Info(viper.Get("data.data.shared_secret"))
	// env
	logger.Sugar().Info(viper.Get("ELB_EXTERNAP_IP"))
	return &ViperConfig{}, nil
}
