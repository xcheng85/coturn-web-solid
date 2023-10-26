package webrtc

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/xcheng85/coturn-web-solid/internal/config"
	"github.com/xcheng85/coturn-web-solid/internal/module"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/handler"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/rest"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/service"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type WebRTCModule struct{}

func (m WebRTCModule) Startup(ctx context.Context, mono module.IModuleContext) error {
	container := dig.New()
	container.Provide(func() *zap.Logger {
		return mono.Logger()
	})
	container.Provide(func() config.IConfig {
		return mono.Config()
	})
	container.Provide(func() *chi.Mux {
		return mono.Mux()
	})
	container.Provide(handler.NewWebRTCHandler)
	container.Provide(rest.NewWebRTCRouter)
	container.Provide(service.NewWebRTCService)
	container.Provide(func() context.Context {
		return ctx
	})
	err := container.Invoke(func(r *rest.WebRTCRouter) error {
		return r.Register()
	})
	return err
}
func NewWebRTCModule() module.Module {
	return &WebRTCModule{}
}
