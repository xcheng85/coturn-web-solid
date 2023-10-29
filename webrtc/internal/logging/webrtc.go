package logging

import (
	"context"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/domain"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/dto"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/service"
	"go.uber.org/zap"
)

// decorator pattern in golang, automatically satisfy interface: handler.RendererHandler
type WebRTCService struct {
	service.IWebRTCService
	logger *zap.Logger
}

var _ service.IWebRTCService = (*WebRTCService)(nil)

// ctor
func LogServiceAccess(service service.IWebRTCService, logger* zap.Logger) WebRTCService {
	return WebRTCService{
		service,
		logger,
	}
}

func (svc WebRTCService) GetWebRTCConfig(ctx context.Context, data dto.GetWebRTCConfigRequest) (rtcconfig *domain.RTCConfig, err error) {
	svc.logger.Info("--> webrtc.GetWebRTCConfig")
	defer func() {
		svc.logger.Sugar().Error(err)
		svc.logger.Info("<-- webrtc.GetWebRTCConfig")
	}()
	return svc.IWebRTCService.GetWebRTCConfig(ctx, data)
}
