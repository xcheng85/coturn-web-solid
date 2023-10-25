package module

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/xcheng85/coturn-web-solid/internal/config"
	"go.uber.org/zap"
)

// chi.Mux is the implementation of chi.Router interface
type IModuleContext interface {
	Mux() *chi.Mux
	Logger() *zap.Logger
	Config() config.IConfig
}

type Module interface {
	Startup(context.Context, IModuleContext) error
}