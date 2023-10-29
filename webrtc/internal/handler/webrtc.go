package handler

import (
	"context"
	"net/http"
	"github.com/go-chi/render"
	"github.com/xcheng85/coturn-web-solid/internal/auth"
	http_utils "github.com/xcheng85/coturn-web-solid/internal/http"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/dto"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/service"
)

type IWebRTCHandler interface {
	GetWebRTCConfig(w http.ResponseWriter, r *http.Request)
	// middleware
	Authorize(next http.Handler) http.Handler
}

type webRTCHandler struct {
	webRTCService service.IWebRTCService
	authService auth.IAuthService
}

func NewWebRTCHandler(webRTCService service.IWebRTCService, authService auth.IAuthService) IWebRTCHandler {
	return &webRTCHandler{
		webRTCService,
		authService,
	}
}

func (handler webRTCHandler) GetWebRTCConfig(w http.ResponseWriter, r *http.Request) {
	data := &dto.GetWebRTCConfigRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, http_utils.ErrBadRequest(err))
		return
	}
	ctx := r.Context()
	webrtcConfig, err := handler.webRTCService.GetWebRTCConfig(ctx, *data)
	if err != nil {
		render.Render(w, r, http_utils.ErrServerInternal(err))
		return
	}
	// convert domain to dto response
	render.Status(r, http.StatusOK)
	render.Render(w, r, &dto.Response{RTCConfig: webrtcConfig})
}

func (handler webRTCHandler) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := handler.authService.Authorize(r)
		if err != nil {
			render.Render(w, r, http_utils.ErrUnauthorized(err))
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
