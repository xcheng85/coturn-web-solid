package rest

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/handler"
)

type WebRTCRouter struct {
	handler handler.WebRTCHandler
	mux *chi.Mux
	ctx context.Context
}

func NewWebRTCRouter(handler handler.WebRTCHandler, ctx context.Context, mux *chi.Mux) *WebRTCRouter {
	return &WebRTCRouter{
		handler: handler,
		mux: mux,
		ctx: ctx,
	}
}

func (router* WebRTCRouter) Register() error {
	r := chi.NewRouter()
	r.Use(router.handler.Authorize)
	r.Get("/", router.handler.GetWebRTCConfig)
	router.mux.Mount("/", r)
	return nil
}
