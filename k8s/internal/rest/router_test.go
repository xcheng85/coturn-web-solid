package rest

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/xcheng85/coturn-web-solid/k8s/internal/handler"
	"github.com/xcheng85/coturn-web-solid/internal/test"
)

func TestNewK8sRouter_RegisterLivenessProbe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	mux := chi.NewRouter()
	mockK8sHandler := &handler.MockIK8sHandler{}
	mockK8sHandler.On("GetLivenessProbe", mock.Anything, mock.Anything).Return().Once()
	k8sHandler := NewK8sRouter(mockK8sHandler, ctx, mux)
	k8sHandler.Register()
	ts := httptest.NewServer(mux)
	defer ts.Close()
	// GET /k8s/livenessProbe
	if _, body := test.TestRequest(t, ts, "GET", "/k8s/livenessProbe", nil); body != "" {
		t.Fatalf(body)
	}
}
