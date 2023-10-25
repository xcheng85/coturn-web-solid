package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(cfg LogConfig) (logger *zap.Logger) {
	switch cfg.Environment {
	case "production":
		logger, _ = zap.NewProductionConfig().Build()
	default:
		logger, _ = zap.NewDevelopmentConfig().Build()
	}
	return logger
}

// avoid naming conflict
func logLevelToZap(level Level) zapcore.Level {
	switch level {
	case PANIC:
		return zapcore.PanicLevel
	case ERROR:
		return zapcore.ErrorLevel
	case WARN:
		return zapcore.WarnLevel
	case INFO:
		return zapcore.InfoLevel
	case DEBUG:
		return zapcore.DebugLevel
	case FATAL:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
