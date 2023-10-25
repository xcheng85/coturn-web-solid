package main

import (
	"fmt"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-chi/chi/v5/middleware"
	"github.com/xcheng85/coturn-web-solid/internal/config"
	"github.com/xcheng85/coturn-web-solid/internal/logger"
	"github.com/xcheng85/coturn-web-solid/internal/module"
	"github.com/xcheng85/coturn-web-solid/internal/worker"
	"github.com/xcheng85/coturn-web-solid/k8s"
	"github.com/xcheng85/coturn-web-solid/webrtc"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func main() {
	container := dig.New()
	container.Provide(newContext)
	err := container.Provide(
		func() *zap.Logger {
			return logger.NewZapLogger(logger.LogConfig{
				LogLevel: logger.DEBUG,
			})
		})
	if err != nil {
		panic(err)
	}
	container.Provide(
		func(logger *zap.Logger) (config.IConfig, error) {
			return config.NewViperConfig([]string{os.Getenv("CONFIG_PATH"), os.Getenv("SECRET_PATH")}, logger)
		})
	container.Provide(newCompositionRoot)
	container.Provide(k8s.NewK8sModule, dig.Name("k8s"))
	container.Provide(webrtc.NewWebRTCModule, dig.Name("webrtc"))
	container.Provide(newMux)
	container.Provide(newModuleContext)
	container.Provide(worker.NewWorkerSyncer)

	//err = container.Invoke(func(p module.IModuleContext) error {
	// err = container.Invoke(func(p worker.IWorkerSyncer) error {
	// 	return nil
	// })

	// err = container.Invoke(func(p struct {
	// 	dig.In
	// 	K8s    module.Module `name:"k8s"`
	// 	WebRTC module.Module `name:"webrtc"`
	// 	Mux    *chi.Mux
	// 	WorkerSyncer worker.IWorkerSyncer
	// }) error {
	// 	fmt.Println(p.Mux)
	// 	newCompositionRoot(p.K8s, p.WebRTC, p.WorkerSyncer)
	// 	return nil
	// })
	err = container.Invoke(func(p struct {
		dig.In
		ModuleContext module.IModuleContext
		K8s           module.Module `name:"k8s"`
		WebRTC        module.Module `name:"webrtc"`
		Mux           *chi.Mux
		WorkerSyncer  worker.IWorkerSyncer
	}) error {
		root := newCompositionRoot(p.ModuleContext, p.K8s, p.WebRTC, p.WorkerSyncer)
		root.startupModules()
		return nil
	})

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
