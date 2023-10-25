package webrtc

import (
	"context"
	"github.com/xcheng85/coturn-web-solid/internal/module"
)

type WebRTCModule struct{}

func (m WebRTCModule) Startup(ctx context.Context, mono module.IModuleContext) error {
	return nil
}
func NewWebRTCModule() module.Module {
	return &WebRTCModule{}
}
