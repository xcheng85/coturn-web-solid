package dto

import (
	"net/http"
	"github.com/xcheng85/coturn-web-solid/webrtc/internal/domain"
)

type (
	Response struct {
		*domain.RTCConfig
	}
)

// interface to satisfiy for go-chi render
func (rd *Response) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}