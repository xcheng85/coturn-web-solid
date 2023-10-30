package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xcheng85/coturn-web-solid/internal/auth"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/domain"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/dto"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/service"
)

func TestWebRTCHandler_GetWebRTCConfig(t *testing.T) {
	scenarios := []struct {
		desc                string
		inRequest           func() *http.Request
		inWebRTCServiceMock func() *service.MockIWebRTCService
		inAuthServiceMock   func() *auth.MockIAuthService
		expectedStatus      int
		expectedPayload     string
	}{
		{
			desc: "Happy Path",
			inRequest: func() *http.Request {
				// perform request
				requestData := &dto.GetWebRTCConfigRequest{}
				data, _ := json.Marshal(requestData)
				body := bytes.NewBuffer(data)
				request, err := http.NewRequest("GET", "/", body)
				// define context and therefore test timeout
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				request = request.WithContext(ctx)
				request.Header.Set("Content-Type", "application/json")
				require.NoError(t, err)
				return request
			},
			inWebRTCServiceMock: func() *service.MockIWebRTCService {
				output := &domain.RTCConfig{
					LifetimeDuration:   "3600s",
					BlockStatus:        "NOT_BLOCKED",
					IceTransportPolicy: "all",
					IceServers:         []domain.ICEServer{},
				}

				mockWebRTCService := &service.MockIWebRTCService{}
				mockWebRTCService.On("GetWebRTCConfig", mock.Anything, mock.Anything).Return(output, nil).Once()
				return mockWebRTCService
			},
			inAuthServiceMock: func() *auth.MockIAuthService {
				mockAuthService := &auth.MockIAuthService{}
				return mockAuthService
			},
			expectedStatus:  http.StatusOK,
			expectedPayload: "{\"lifetimeDuration\":\"3600s\",\"iceServers\":[],\"blockStatus\":\"NOT_BLOCKED\",\"iceTransportPolicy\":\"all\"}\n",
		},
		{
			desc: "Server error: No External Ips",
			inRequest: func() *http.Request {
				// perform request
				requestData := &dto.GetWebRTCConfigRequest{}
				data, _ := json.Marshal(requestData)
				body := bytes.NewBuffer(data)
				request, err := http.NewRequest("GET", "/", body)
				// define context and therefore test timeout
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				request = request.WithContext(ctx)
				request.Header.Set("Content-Type", "application/json")
				require.NoError(t, err)
				return request
			},
			inWebRTCServiceMock: func() *service.MockIWebRTCService {
				mockWebRTCService := &service.MockIWebRTCService{}
				mockWebRTCService.On("GetWebRTCConfig", mock.Anything, mock.Anything).Return(nil, service.NewEmptyExternalIpErr()).Once()
				return mockWebRTCService
			},
			inAuthServiceMock: func() *auth.MockIAuthService {
				mockAuthService := &auth.MockIAuthService{}
				return mockAuthService
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedPayload: "{\"status\":\"Server Internal Error\",\"error\":\"no external ips of load balancer(s) are available\"}\n",
		},
	}
	for _, s := range scenarios {
		scenario := s
		t.Run(scenario.desc, func(t *testing.T) {
			mockWebRTCService := scenario.inWebRTCServiceMock()
			mockAuthService := scenario.inAuthServiceMock()
			// build handler
			webrtcHandler := NewWebRTCHandler(mockWebRTCService, mockAuthService)
			// perform request
			response := httptest.NewRecorder()
			webrtcHandler.GetWebRTCConfig(response, scenario.inRequest())

			// validate outputs
			require.Equal(t, scenario.expectedStatus, response.Code, scenario.desc)

			payload, _ := io.ReadAll(response.Body)
			assert.Equal(t, scenario.expectedPayload, string(payload), scenario.desc)
		})
	}
}

type MockHttpHandler struct {
	mock.Mock
}

// DoSomething is a method on MyMockedObject that implements some interface
func (m *MockHttpHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	return
}

func TestWebRTCHandler_Authorize(t *testing.T) {
	scenarios := []struct {
		desc                string
		inRequest           func() *http.Request
		inWebRTCServiceMock func() *service.MockIWebRTCService
		inAuthServiceMock   func() *auth.MockIAuthService
		inHttpHandlerMock   func() http.Handler
		expectedStatus      int
		expectedPayload     string
	}{
		{
			desc: "Happy Path",
			inRequest: func() *http.Request {
				// perform request
				requestData := &dto.GetWebRTCConfigRequest{}
				data, _ := json.Marshal(requestData)
				body := bytes.NewBuffer(data)
				request, err := http.NewRequest("GET", "/", body)
				// define context and therefore test timeout
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				request = request.WithContext(ctx)
				request.Header.Set("Content-Type", "application/json")
				require.NoError(t, err)
				return request
			},
			inWebRTCServiceMock: func() *service.MockIWebRTCService {
				output := &domain.RTCConfig{
					LifetimeDuration:   "3600s",
					BlockStatus:        "NOT_BLOCKED",
					IceTransportPolicy: "all",
					IceServers:         []domain.ICEServer{},
				}
				mockWebRTCService := &service.MockIWebRTCService{}
				mockWebRTCService.On("GetWebRTCConfig", mock.Anything, mock.Anything).Return(output, nil).Once()
				return mockWebRTCService
			},
			inAuthServiceMock: func() *auth.MockIAuthService {
				mockAuthService := &auth.MockIAuthService{}
				mockAuthService.On("Authorize", mock.Anything).Return("user", nil).Once()
				return mockAuthService
			},
			inHttpHandlerMock: func() http.Handler {
				mockHttpHandler := new(MockHttpHandler)
				mockHttpHandler.On("ServeHTTP", mock.Anything, mock.Anything).Return()
				return mockHttpHandler
			},
			expectedStatus:  http.StatusOK,
			expectedPayload: "",
		},
		{
			desc: "401: Unauthorized",
			inRequest: func() *http.Request {
				// perform request
				requestData := &dto.GetWebRTCConfigRequest{}
				data, _ := json.Marshal(requestData)
				body := bytes.NewBuffer(data)
				request, err := http.NewRequest("GET", "/", body)
				// define context and therefore test timeout
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				request = request.WithContext(ctx)
				request.Header.Set("Content-Type", "application/json")
				require.NoError(t, err)
				return request
			},
			inWebRTCServiceMock: func() *service.MockIWebRTCService {
				output := &domain.RTCConfig{
					LifetimeDuration:   "3600s",
					BlockStatus:        "NOT_BLOCKED",
					IceTransportPolicy: "all",
					IceServers:         []domain.ICEServer{},
				}
				mockWebRTCService := &service.MockIWebRTCService{}
				mockWebRTCService.On("GetWebRTCConfig", mock.Anything, mock.Anything).Return(output, nil).Once()
				return mockWebRTCService
			},
			inAuthServiceMock: func() *auth.MockIAuthService {
				mockAuthService := &auth.MockIAuthService{}
				mockAuthService.On("Authorize", mock.Anything).Return("", auth.NewUnauthorizedError("empty bearerToken")).Once()
				return mockAuthService
			},
			inHttpHandlerMock: func() http.Handler {
				mockHttpHandler := new(MockHttpHandler)
				mockHttpHandler.On("ServeHTTP", mock.Anything, mock.Anything).Return()
				return mockHttpHandler
			},
			expectedStatus:  http.StatusUnauthorized,
			expectedPayload: "{\"status\":\"Not Authorized\",\"error\":\"empty bearerToken\"}\n",
		},
	}
	for _, s := range scenarios {
		scenario := s
		t.Run(scenario.desc, func(t *testing.T) {
			mockWebRTCService := scenario.inWebRTCServiceMock()
			mockAuthService := scenario.inAuthServiceMock()
			mockHttpHandler := scenario.inHttpHandlerMock()
			// build handler
			webrtcHandler := NewWebRTCHandler(mockWebRTCService, mockAuthService)

			// perform request
			response := httptest.NewRecorder()
			handler := webrtcHandler.Authorize(mockHttpHandler)
			handler.ServeHTTP(response, scenario.inRequest())

			// validate outputs
			require.Equal(t, scenario.expectedStatus, response.Code, scenario.desc)

			payload, _ := io.ReadAll(response.Body)
			assert.Equal(t, scenario.expectedPayload, string(payload), scenario.desc)
		})
	}
}
