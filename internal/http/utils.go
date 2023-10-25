package http

import (
	"github.com/go-chi/render"
	"net/http"
)

type HttpResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code
	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *HttpResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

var ErrNotFound = &HttpResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

func ErrBadRequest(err error) render.Renderer {
	return &HttpResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Bad request",
		ErrorText:      err.Error(),
	}
}

func ErrServerInternal(err error) render.Renderer {
	return &HttpResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Server Internal Error",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &HttpResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering dto response",
		ErrorText:      err.Error(),
	}
}

func TextOkRender(message string) render.Renderer {
	return &HttpResponse{
		Err:            nil,
		HTTPStatusCode: 200,
		StatusText:     message,
		ErrorText:      "",
	}
}