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
	service.WebRTCService
	logger *zap.Logger
}

var _ service.WebRTCService = (*WebRTCService)(nil)

// ctor
func LogServiceAccess(service service.WebRTCService, logger* zap.Logger) WebRTCService {
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
	return svc.WebRTCService.GetWebRTCConfig(ctx, data)
}
