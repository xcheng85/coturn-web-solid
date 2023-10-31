package rest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	http_utils "github.com/xcheng85/coturn-web-solid/internal/http"
	"github.com/xcheng85/coturn-web-solid/internal/test"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/handler"
)

func TestNewK8sRouter_Register(t *testing.T) {
	scenarios := []struct {
		desc                string
		mux                 *chi.Mux
		inMockWebRTCHandler func() *handler.MockIWebRTCHandler
		expectedStatus      int
		expectedPayload     string
	}{
		{
			desc: "Happy Path",
			mux:  chi.NewRouter(),
			inMockWebRTCHandler: func() *handler.MockIWebRTCHandler {
				mockWebRTCHandler := &handler.MockIWebRTCHandler{}
				mockWebRTCHandler.On("GetWebRTCConfig", mock.Anything, mock.Anything).Return(nil).Once()
				mockWebRTCHandler.On("Authorize", mock.Anything).Return(new(test.MockHttpHandler)).Once()
				return mockWebRTCHandler
			},
			expectedStatus:  http.StatusOK,
			expectedPayload: "",
		},
		{
			desc: "Not Authorized",
			mux:  chi.NewRouter(),
			inMockWebRTCHandler: func() *handler.MockIWebRTCHandler {
				mockWebRTCHandler := &handler.MockIWebRTCHandler{}
				mockWebRTCHandler.On("GetWebRTCConfig", mock.Anything, mock.Anything).Return(nil).Once()
				f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					render.Render(w, r, http_utils.ErrUnauthorized(errors.New("error")))
					return
				})
				mockWebRTCHandler.On("Authorize", mock.Anything).Return(f).Once()
				return mockWebRTCHandler
			},
			expectedStatus:  http.StatusUnauthorized,
			expectedPayload: "{\"status\":\"Not Authorized\",\"error\":\"error\"}\n",
		},
	}

	for _, s := range scenarios {
		scenario := s
		t.Run(scenario.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			mux := scenario.mux
			mockWebRTCHandler := scenario.inMockWebRTCHandler()
			webRTCHandler := NewWebRTCRouter(mockWebRTCHandler, ctx, mux)
			webRTCHandler.Register()
			ts := httptest.NewServer(mux)
			defer ts.Close()

			response, body := test.TestRequest(t, ts, "GET", "/", nil)
			require.Equal(t, scenario.expectedStatus, response.StatusCode, scenario.desc)
			assert.Equal(t, scenario.expectedPayload, body, scenario.desc)
		})
	}
}
