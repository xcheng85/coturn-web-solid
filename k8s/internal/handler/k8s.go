package handler

import (
	"github.com/go-chi/render"
	http_utils "github.com/xcheng85/coturn-web-solid/internal/http"
	"net/http"
)

//go:generate mockery --name IK8sHandler
type IK8sHandler interface {
	GetLivenessProbe(w http.ResponseWriter, r *http.Request)
	GetReadinessProbe(w http.ResponseWriter, r *http.Request)
}

type k8sHandler struct {
}

func NewK8sHandler() IK8sHandler {
	return &k8sHandler{}
}

func (handler k8sHandler) GetLivenessProbe(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.Render(w, r, http_utils.TextOkRender("livenessProbe passes"))
}

func (handler k8sHandler) GetReadinessProbe(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.Render(w, r, http_utils.TextOkRender("readinessProbe passes"))
}
