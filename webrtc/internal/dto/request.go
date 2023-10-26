package dto

import (
	"net/http"
)

type (
	GetWebRTCConfigRequest struct {
	}
)

func (a *GetWebRTCConfigRequest) Bind(r *http.Request) error {
	return nil
}
