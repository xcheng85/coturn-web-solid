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
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/dto"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/service"
)

func TestWebRTCHandler_GetWebRTCConfig_NoExternalIps(t *testing.T) {
	mockWebRTCService := &service.MockIWebRTCService{}
	mockAuthService := &auth.MockIAuthService{}

	// GetWebRTCConfig has two parameters
	mockWebRTCService.On("GetWebRTCConfig", mock.Anything, mock.Anything).Return(nil, service.NewEmptyExternalIpErr()).Once()
	webrtcHandler := NewWebRTCHandler(mockWebRTCService, mockAuthService)
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
	response := httptest.NewRecorder()
	webrtcHandler.GetWebRTCConfig(response, request)
	require.Equal(t, 500, response.Code)
	payload, _ := io.ReadAll(response.Body)
	assert.Equal(t, "{\"status\":\"Server Internal Error\",\"error\":\"no external ips of load balancer(s) are available\"}\n",
		string(payload), "Server error: No External Ips")
}
