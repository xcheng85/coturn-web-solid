package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/xcheng85/coturn-web-solid/internal/auth"
	"github.com/xcheng85/coturn-web-solid/internal/config"
	"github.com/xcheng85/coturn-web-solid/internal/module"
	"github.com/xcheng85/coturn-web-solid/internal/worker"
	_ "go.uber.org/dig"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// composition root
// application in the hexongal arch
// app must implement module interface, which is required in each sub module
// owner of all modules
type CompositionRoot struct {
	moduleCtx    module.IModuleContext
	modules      []module.Module
	workerSyncer worker.IWorkerSyncer
}

func newCompositionRoot(moduleCtx module.IModuleContext, k8s, webrtc module.Module, workerSyncer worker.IWorkerSyncer) *CompositionRoot {
	return &CompositionRoot{
		moduleCtx:    moduleCtx,
		modules:      []module.Module{k8s, webrtc},
		workerSyncer: workerSyncer,
	}
}

func (r *CompositionRoot) startupModules() error {
	for _, module := range r.modules {
		if err := module.Startup(r.workerSyncer.Context(), r.moduleCtx); err != nil {
			return err
		}
	}

	r.workerSyncer.Add(r.runRestServer)
	return r.workerSyncer.Sync()
}

// worker for running Rest server for reverse proxy
func (r *CompositionRoot) runRestServer(ctx context.Context) error {
	mux := r.moduleCtx.Mux()
	logger := r.moduleCtx.Logger()
	config := r.moduleCtx.Config()
	address := fmt.Sprintf(":%d", config.Get("port"))
	logger.Sugar().Infof("runRestServer: %d", address)
	restServer := &http.Server{
		Addr: address,
		Handler: mux,
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		logger.Sugar().Info("web server started")
		defer logger.Sugar().Info("web server shutdown")
		if err := restServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	group.Go(func() error {
		// received cancel signal from the derived
		<-gCtx.Done()
		logger.Sugar().Info("web server to be shutdown")
		// gracefully shut down rest server
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		if err := restServer.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})
	// block here
	return group.Wait()
}

func newMux() *chi.Mux {
	mux := chi.NewRouter()
	// A good base middleware stack
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	// mux.Use(middleware.Timeout(1 * time.Second))
	mux.Use(middleware.URLFormat)
	mux.Use(render.SetContentType(render.ContentTypeJSON))
	return mux
}

func newContext() context.Context {
	return context.Background()
}

type ModuleContext struct {
	mux    *chi.Mux
	logger *zap.Logger
	config config.IConfig
	auth   auth.IAuthService
}

func newModuleContext(mux *chi.Mux, logger *zap.Logger, config config.IConfig, auth auth.IAuthService) module.IModuleContext {
	return &ModuleContext{
		mux,
		logger,
		config,
		auth,
	}
}

func (r *ModuleContext) Mux() *chi.Mux {
	return r.mux
}

func (r *ModuleContext) Logger() *zap.Logger {
	return r.logger
}

func (r *ModuleContext) Config() config.IConfig {
	return r.config
}

func (r *ModuleContext) Auth() auth.IAuthService {
	return r.auth
}
