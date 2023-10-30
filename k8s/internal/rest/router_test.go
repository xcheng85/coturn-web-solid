package rest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/xcheng85/coturn-web-solid/k8s/internal/handler"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}

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
	if _, body := testRequest(t, ts, "GET", "/k8s/livenessProbe", nil); body != "" {
		t.Fatalf(body)
	}
}
