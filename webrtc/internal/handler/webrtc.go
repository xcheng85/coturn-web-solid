package handler

import (
	"github.com/go-chi/render"
	http_utils "github.com/xcheng85/coturn-web-solid/internal/http"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/dto"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/service"
	"net/http"
)

type WebRTCHandler interface {
	GetWebRTCConfig(w http.ResponseWriter, r *http.Request)
}

type webRTCHandler struct {
	service service.WebRTCService
}

func NewWebRTCHandler(service service.WebRTCService) WebRTCHandler {
	return &webRTCHandler{
		service,
	}
}

func (handler webRTCHandler) GetWebRTCConfig(w http.ResponseWriter, r *http.Request) {
	data := &dto.GetWebRTCConfigRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, http_utils.ErrBadRequest(err))
		return
	}
	ctx := r.Context()
	webrtcConfig, err := handler.service.GetWebRTCConfig(ctx, *data)
	if err != nil {
		render.Render(w, r, http_utils.ErrServerInternal(err))
		return
	}
	// convert domain to dto response
	render.Status(r, http.StatusOK)
	render.Render(w, r, &dto.Response{RTCConfig: webrtcConfig})
}
