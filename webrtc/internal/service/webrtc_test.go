package service

import (
	context "context"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"github.com/xcheng85/coturn-web-solid/internal/config"
	"github.com/xcheng85/coturn-web-solid/internal/logger"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/domain"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/dto"
)

func TestWebRTCService_GetWebRTCConfig(t *testing.T) {
	scenarios := []struct {
		desc           string
		inLogger       *zap.Logger
		inConfigMock   func() *config.MockIConfig
		expectedConfig *domain.RTCConfig
		expectedError  error
	}{
		{
			desc: "Happy Path",
			inLogger: logger.NewZapLogger(logger.LogConfig{
				LogLevel: logger.DEBUG,
			}),
			inConfigMock: func() *config.MockIConfig {
				mockConfig := &config.MockIConfig{}
				mockConfig.On("Get", "ELB_EXTERNAP_IP").Return("1.2.3.4", nil).Once()
				mockConfig.On("Get", "turn_config.use_public_stun").Return(true, nil).Once()
				mockConfig.On("Get", "turn_config.use_private_stun").Return(true, nil).Once()
				mockConfig.On("Get", "turn_config.use_tcp").Return(true, nil).Once()
				mockConfig.On("Get", "turn_config.use_udp").Return(true, nil).Once()
				mockConfig.On("Get", "turn_config.port").Return(1234, nil).Once()
				mockConfig.On("Get", "turn_config.username").Return("username", nil).Once()
				mockConfig.On("Get", "turn_config.ttl_seconds").Return(3600, nil).Once()
				mockConfig.On("Get", "turn_config.ice_transport_policy").Return("all", nil).Once()
				mockConfig.On("Get", "data.data.password").Return("password", nil).Once()
				return mockConfig
			},
			expectedConfig: &domain.RTCConfig{
				LifetimeDuration: "3600s",
				IceServers: []domain.ICEServer{
					{URLs: []string{"stun:stun.l.google.com:19302", "stun:1.2.3.4:1234"},
						UserName: "", Credential: ""},
					{URLs: []string{"turn:1.2.3.4:1234?transport=tcp", "turn:1.2.3.4:1234?transport=udp"},
						UserName: "username", Credential: "password"}},
				BlockStatus:        "NOT_BLOCKED",
				IceTransportPolicy: "all",
			},
			expectedError: nil,
		},
		{
			desc: "Empty External Ips",
			inLogger: logger.NewZapLogger(logger.LogConfig{
				LogLevel: logger.DEBUG,
			}),
			inConfigMock: func() *config.MockIConfig {
				mockConfig := &config.MockIConfig{}
				mockConfig.On("Get", "ELB_EXTERNAP_IP").Return("1.2", nil).Once()
				mockConfig.On("Get", "turn_config.use_public_stun").Return(false, nil).Once()
				mockConfig.On("Get", "turn_config.use_private_stun").Return(true, nil).Once()
				mockConfig.On("Get", "turn_config.use_tcp").Return(false, nil).Once()
				mockConfig.On("Get", "turn_config.use_udp").Return(true, nil).Once()
				mockConfig.On("Get", "turn_config.port").Return(1234, nil).Once()
				mockConfig.On("Get", "turn_config.username").Return("username", nil).Once()
				mockConfig.On("Get", "turn_config.ttl_seconds").Return(3600, nil).Once()
				mockConfig.On("Get", "turn_config.ice_transport_policy").Return("all", nil).Once()
				mockConfig.On("Get", "data.data.password").Return("password", nil).Once()
				return mockConfig
			},
			expectedConfig: nil,
			expectedError:  NewInvalidExternalIpErr("1.2"),
		},
	}
	for _, s := range scenarios {
		scenario := s
		t.Run(scenario.desc, func(t *testing.T) {
			logger := scenario.inLogger
			config := scenario.inConfigMock()
			// define context and therefore test timeout
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			webrtcService := NewWebRTCService(logger, config)
			webrtcConfig, err := webrtcService.GetWebRTCConfig(ctx, dto.GetWebRTCConfigRequest{})
			assert.Equal(t, scenario.expectedConfig, webrtcConfig, scenario.desc)
			assert.Equal(t, scenario.expectedError, err, scenario.desc)
		})
	}
}
