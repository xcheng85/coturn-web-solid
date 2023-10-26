package service

import (
	"context"
	"fmt"
	"github.com/xcheng85/coturn-web-solid/internal/config"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/domain"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/dto"
	"go.uber.org/zap"
)

type WebRTCService interface {
	GetWebRTCConfig(ctx context.Context, data dto.GetWebRTCConfigRequest) (*domain.RTCConfig, error)
}

type webRTCService struct {
	logger *zap.Logger
	config config.IConfig
}

var _ WebRTCService = (*webRTCService)(nil)

func NewWebRTCService(logger *zap.Logger, config config.IConfig) WebRTCService {
	return &webRTCService{
		logger,
		config,
	}
}

func (svc webRTCService) GetWebRTCConfig(ctx context.Context, data dto.GetWebRTCConfigRequest) (*domain.RTCConfig, error) {
	// env
	externalIP := svc.config.Get("ELB_EXTERNAP_IP").(string)
	// config
	usePublicStun := svc.config.Get("turn_config.use_public_stun").(bool)
	usePrivateStun := svc.config.Get("turn_config.use_private_stun").(bool)
	useTcp := svc.config.Get("turn_config.use_tcp").(bool)
	useUdp := svc.config.Get("turn_config.use_udp").(bool)
	port := svc.config.Get("turn_config.port").(int)
	username := svc.config.Get("turn_config.username").(string)
	ttlSeconds := svc.config.Get("turn_config.ttl_seconds").(int)
	iceTransportPolicy := svc.config.Get("turn_config.ice_transport_policy").(string)
	// secret
	password := svc.config.Get("data.data.password").(string)

	ips := []string{externalIP}

	if len(ips) == 0 {
		return nil, &EmptyExternalIpErr{}
	}

	publicStunServerUrls, stunServerUrls, turnServerUrls := []string{
		"stun:stun.l.google.com:19302",
	}, []string{}, []string{}

	if usePublicStun {
		stunServerUrls = append(stunServerUrls, publicStunServerUrls...)
	}
	for _, ip := range ips {
		if usePrivateStun {
			stunServerUrls = append(stunServerUrls, fmt.Sprintf("stun:%s:%d", ip, port))
		}
		if useTcp {
			turnServerUrls = append(turnServerUrls, fmt.Sprintf("turn:%s:%d?transport=tcp", ip, port))
		}
		if useUdp {
			turnServerUrls = append(turnServerUrls, fmt.Sprintf("turn:%s:%d?transport=udp", ip, port))
		}
	}
	svc.logger.Sugar().Info("turnServerUrls:", turnServerUrls)

	iceServers := []domain.ICEServer{}
	if len(stunServerUrls) > 0 {
		iceServers = append(iceServers, domain.ICEServer{
			URLs: stunServerUrls,
		})
	}
	if len(turnServerUrls) > 0 {
		iceServers = append(iceServers, domain.ICEServer{
			URLs:       turnServerUrls,
			UserName:   username,
			Credential: password,
		})
	}
	svc.logger.Sugar().Info("iceServers:", iceServers)

	rtcConfig := &domain.RTCConfig{
		LifetimeDuration:   fmt.Sprintf("%ds", ttlSeconds),
		BlockStatus:        "NOT_BLOCKED",
		IceTransportPolicy: iceTransportPolicy,
		IceServers:         iceServers,
	}

	return rtcConfig, nil
}
