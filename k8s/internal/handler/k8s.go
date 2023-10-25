package handler

import (
	"github.com/go-chi/render"
	http_utils "github.com/xcheng85/coturn-web-solid/internal/http"
	"net/http"
)

type K8sHandler interface {
	GetLivenessProbe(w http.ResponseWriter, r *http.Request)
	GetReadinessProbe(w http.ResponseWriter, r *http.Request)
}

type k8sHandler struct {
}

func NewK8sHandler() K8sHandler {
	return &k8sHandler{}
}

func (handler k8sHandler) GetLivenessProbe(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusCreated)
	render.Render(w, r, http_utils.TextOkRender("livenessProbe passes"))
}

func (handler k8sHandler) GetReadinessProbe(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusCreated)
	render.Render(w, r, http_utils.TextOkRender("readinessProbe passes"))
}
